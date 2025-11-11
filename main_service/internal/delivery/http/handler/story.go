package handler

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"server/internal/helper"
	"server/internal/models"
	"server/internal/service/repository"
	package_log "server/pkg/logging"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type StoryHandler struct {
	logger  *package_log.Logger
	service repository.StoryService
}

func NewStoryHandler(logger *package_log.Logger, service repository.StoryService) *StoryHandler {
	return &StoryHandler{
		logger:  logger,
		service: service,
	}
}

func (h *StoryHandler) StoryRegisterRoutes(r *gin.RouterGroup) {
	r.POST("/add-photo", h.StoriesUploadFile)
	r.POST("/add-video", h.StoriesUploadVideo)
	r.GET("/get-user-stories", h.GetUserStories)
	r.GET("/get-all-stories-for-user", h.GetAllStoryForUser)
	r.GET("/get-story-by-id/:uuid", h.GetStoryById)
}

func (h *StoryHandler) StoriesUploadFile(c *gin.Context) {
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")

	userId, err := helper.IntId(c)
	if err != nil {
		return
	}

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		h.logger.Errorln("error:", err)
		c.String(http.StatusBadRequest, "error: %v", err)
		return
	}
	defer file.Close()

	// uploads papkasini tayyorlash
	if _, err := os.Stat("uploads"); os.IsNotExist(err) {
		if err := os.MkdirAll("uploads", 0755); err != nil {
			h.logger.Errorln("error:", err)
			c.String(http.StatusInternalServerError, "error with creating folder: %v", err)
			return
		}
	}

	// Fayl uchun UUID nom yaratish
	ext := filepath.Ext(header.Filename)
	fileName := fmt.Sprintf("%s%s", uuid.New().String(), ext)
	outPath := filepath.Join("uploads", fileName)

	outFile, err := os.Create(outPath)
	if err != nil {
		h.logger.Errorln("error:", err)
		c.String(http.StatusInternalServerError, "error: %v", err)
		return
	}
	defer outFile.Close()

	// Progress hisoblash
	buffer := make([]byte, 1024*1024) // 1MB buffer
	var uploaded int64
	fileSize := header.Size

	for {
		n, err := file.Read(buffer)
		if n > 0 {
			_, werr := outFile.Write(buffer[:n])
			if werr != nil {
				c.String(http.StatusInternalServerError, "error: %v", werr)
				return
			}
			uploaded += int64(n)

			percent := float64(uploaded) / float64(fileSize) * 100
			fmt.Fprintf(c.Writer, "data: %.2f%%\n\n", percent)
			c.Writer.Flush()
		}

		if err != nil {
			if err == io.EOF {
				break
			}
			h.logger.Errorln("error:", err)
			c.String(http.StatusInternalServerError, "error: %v", err)

			return
		}
	}
	err = h.service.SaveStory(c.Request.Context(), userId, "/uploads/"+fileName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})

		return
	}

	// Yakuniy javob
	fmt.Fprintf(c.Writer, "data: Successfully uploaded! File name: %s\n\n", fileName)
	c.Writer.Flush()
}

func (h *StoryHandler) StoriesUploadVideo(c *gin.Context) {
	// SSE header'lar
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")

	userId, err := helper.IntId(c)
	if err != nil {
		return
	}

	// Faylni olish
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		h.logger.Errorln("error:", err)
		fmt.Fprintf(c.Writer, "data: error: %v\n\n", err)
		c.Writer.Flush()
		return
	}
	defer file.Close()

	// Faqat .mp4 faylga ruxsat
	ext := strings.ToLower(filepath.Ext(header.Filename))
	if ext != ".mp4" {
		fmt.Fprintf(c.Writer, "data: error: faqat .mp4 formatiga ruxsat beriladi\n\n")
		c.Writer.Flush()
		return
	}

	// Fayl uchun papka yaratish
	uid := uuid.New().String()
	dirPath := filepath.Join("uploads/videos", uid)
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		h.logger.Errorln("error:", err)
		fmt.Fprintf(c.Writer, "data: error: katalog yaratishda xato\n\n")
		c.Writer.Flush()
		return
	}

	// Fayl yo‘lini tayyorlash
	videoPath := filepath.Join(dirPath, header.Filename)
	outFile, err := os.Create(videoPath)
	if err != nil {
		h.logger.Errorln("error:", err)
		fmt.Fprintf(c.Writer, "data: error: fayl yaratishda xato\n\n")
		c.Writer.Flush()
		return
	}
	defer outFile.Close()

	// Progress hisoblash
	buffer := make([]byte, 1024*1024) // 1MB buffer
	var uploaded int64
	fileSize := header.Size

	for {
		n, err := file.Read(buffer)
		if n > 0 {
			if _, werr := outFile.Write(buffer[:n]); werr != nil {
				h.logger.Errorln("error:", werr)
				fmt.Fprintf(c.Writer, "data: error: yozishda xato: %v\n\n", werr)
				c.Writer.Flush()
				return
			}
			uploaded += int64(n)
			percent := float64(uploaded) / float64(fileSize) * 100

			// Har 1% da progress yuborish
			fmt.Fprintf(c.Writer, "data: %.2f%%\n\n", percent)
			c.Writer.Flush()
		}

		if err != nil {
			if err == io.EOF {
				break
			}
			h.logger.Errorln("error:", err)
			fmt.Fprintf(c.Writer, "data: error: %v\n\n", err)
			c.Writer.Flush()
			return
		}
	}

	// Yuklanish tugadi, endi HLS formatga o‘tkazamiz
	fmt.Fprintf(c.Writer, "File uploaded, HLS convertation starting...\n\n")
	c.Writer.Flush()

	hlsPath := filepath.Join(dirPath, "index.m3u8")
	cmd := exec.Command(
		"ffmpeg", "-i", videoPath,
		"-codec:", "copy",
		"-start_number", "0",
		"-hls_time", "10",
		"-hls_list_size", "0",
		"-f", "hls",
		hlsPath,
	)

	if err := cmd.Run(); err != nil {
		h.logger.Errorln("error:", err)
		fmt.Fprintf(c.Writer, "data: error: HLS konvertatsiyada xato: %v\n\n", err)
		c.Writer.Flush()
		return
	}

	filePath := fmt.Sprintf("/uploads/videos/%s/index.m3u8", uid)

	err = h.service.SaveStory(c.Request.Context(), userId, filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})

		return
	}

	// Yakuniy javob
	fmt.Fprintf(c.Writer, "data: Successfully uploaded! HLS URL: %s", filePath)
	c.Writer.Flush()
}

func (h *StoryHandler) GetUserStories(c *gin.Context) {
	userId, err := helper.IntId(c)
	if err != nil {
		return
	}

	data, err := h.service.GetUserStory(c.Request.Context(), userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	data.User_id = userId

	c.JSON(http.StatusOK, data)
}

func (h *StoryHandler) GetAllStoryForUser(c *gin.Context) {
	userId, err := helper.IntId(c)
	if err != nil {
		return
	}

	limit, _ := strconv.Atoi(c.Query("limit"))
	offset, _ := strconv.Atoi(c.Query("offset"))

	if limit == 0 {
		limit = 5
		offset = 0
	}

	dto := models.GetStoryDto{
		User_Id: userId,
		Limit:   limit,
		Offset:  offset,
	}

	data, err := h.service.GetAllStoryForUser(c.Request.Context(), dto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, data)
}

func (h *StoryHandler) GetStoryById(c *gin.Context) {
	uuid := c.Param("uuid")

	data, err := h.service.GetStoryById(c.Request.Context(), uuid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, data)

}
