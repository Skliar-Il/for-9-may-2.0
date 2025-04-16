package dto

type MedalDTO struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	ImageUrl string `json:"photo_link"`
}

type CreateMedalDTO struct {
	Name      string `json:"name"`
	ImageLink string `json:"photo_link"`
}
