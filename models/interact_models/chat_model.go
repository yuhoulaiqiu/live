package interact_models

import "live/commen/models"

type ChatModel struct {
	models.Model
	RoomNumber string `json:"roomNumber"`
	SendUserId uint   `json:"sendUserId"`
	Msg        string `json:"msg"`
}
