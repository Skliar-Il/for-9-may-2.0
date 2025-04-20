package dto

import "github.com/google/uuid"

type CreatePhotoDTO struct {
	PersonID   uuid.UUID `json:"person_id"`
	Link       string    `json:"link"`
	MainStatus bool      `json:"main_status"`
}

type PhotoDTO struct {
	ID     int    `json:"id"`
	Link   string `json:"link"`
	IsMain bool   `json:"is_main"`
}
