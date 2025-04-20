package dto

import "github.com/google/uuid"

type CreateNewPhotoDTO struct {
	PersonID   uuid.UUID `json:"person_id"`
	Link       string    `json:"link"`
	MainStatus bool      `json:"main_status"`
}

type DeletePhotoDTO struct {
	ID       int       `json:"id"`
	PersonID uuid.UUID `json:"person_id"`
}

type PhotoDTO struct {
	ID   int    `json:"id"`
	Link string `json:"link"`
}
