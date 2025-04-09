package model

type CreatePersonModel struct {
	Name              string `json:"name" binding:"required"`
	Surname           string `json:"surname" binding:"required"`
	Patronymic        string `json:"patronymic" binding:"required"`
	DateBirth         int    `json:"date_birth"`
	DateDeath         int    `json:"date_death"`
	City              string `json:"city"`
	History           string `json:"history" binding:"required"`
	Rank              string `json:"rank" binding:"required"`
	Role              bool   `json:"role" binding:"required"`
	ContactEmail      string `json:"contact_email" binding:"required"`
	ContactName       string `json:"contact_name" binding:"required"`
	ContactSurname    string `json:"contact_surname" binding:"required"`
	ContactPatronymic string `json:"contact_patronymic" binding:"required"`
	ContactTelegram   string `json:"contact_telegram" binding:"required"`
	Medals            []int  `json:"medals" binding:"required"`
	Relative          string `json:"relative" binding:"required"`
}
