package models


type Story struct {
	User_id   string `json:"user_id" db:"uuid"`
	Username  string `json:"username" db:"username"`
	Image_url string `json:"image_url" db:"image_url"`
	File_url  string `json:"file_url" db:"file_url"`
}

type MyStory struct {
	Username  string `json:"username" db:"username"`
	Image_url string `json:"image_url" db:"image_url"`
	File_url  string `json:"file_url" db:"file_url"`
}

type GetStoryDto struct {
	User_Id string
	Limit   int
	Offset  int
}
