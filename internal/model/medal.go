package model

type MedalModel struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	ImageUrl string `json:"image_url"`
}

type CreateMedalModel struct {
	Name string `json:"name"`
}
