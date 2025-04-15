package model

type MedalModel struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	ImageUrl string `json:"photo_link"`
}

type CreateMedalModel struct {
	Name      string `json:"name"`
	ImageLink string `json:"photo_link"`
}
