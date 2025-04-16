package dto

import "github.com/google/uuid"

type CreateNewPhotoDTO struct {
	PersonID   uuid.UUID `json:"person_id"`
	Link       string    `json:"link"`
	MainStatus bool      `json:"main_status"`
}
