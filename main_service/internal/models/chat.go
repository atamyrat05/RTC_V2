package models

import "time"

type Chat struct {
	User_Id   string    `json:"user_id"`
	Username  string    `json:"username"`
	Image_url string    `json:"image_url"`
	LastTime  time.Time `json:"last_time"`
	RoomId    string    `json:"room_id"`
}

type ChatWithUser struct {
	UserID    string
	PartnerID string
	Limit     int
	Offset    int
}

type ChatsMessage struct {
	Message string    `json:"message"`
	Time    time.Time `json:"time"`
	Sender  string    `json:"sender"`
}

type SaveChat struct {
	Message      string
	From_user_id string
	To_user_id   string
}
