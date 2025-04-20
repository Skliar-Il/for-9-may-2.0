package dto

import (
	"github.com/google/uuid"
	"time"
)

type CreateGalleryPostDTO struct {
	Description string    `json:"description"`
	Date        time.Time `json:"date"`
}

type GalleryPostDTO struct {
	ID          uuid.UUID `json:"id"`
	Link        string    `json:"link"`
	Description string    `json:"description"`
	Date        time.Time `json:"date"`
}
