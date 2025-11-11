package models

type CreateRoomReq struct {
	ID string `json:"id" binding:"required"`
}

type RoomRes struct {
	ID   int `json:"id"`
	Name string `json:"name"`
}

type ClientRes struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

type SingleRoom struct {
	User1_id string
	User2_id string
}

type SaveChat struct {
	Message      string
	From_user_id string
	To_user_id   string
}
