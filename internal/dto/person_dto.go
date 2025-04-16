package dto

import "time"

type CreatePersonDTO struct {
	Name              string `json:"name" binding:"required"`
	Surname           string `json:"surname" binding:"required"`
	Patronymic        string `json:"patronymic"`
	DateBirth         int    `json:"date_birth"`
	DateDeath         int    `json:"date_death"`
	City              string `json:"city"`
	History           string `json:"history" binding:"required"`
	Rank              string `json:"rank" binding:"required"`
	Role              bool   `json:"role" binding:"required"`
	ContactEmail      string `json:"contact_email" binding:"required"`
	ContactName       string `json:"contact_name" binding:"required"`
	ContactSurname    string `json:"contact_surname" binding:"required"`
	ContactPatronymic string `json:"contact_patronymic"`
	ContactTelegram   string `json:"contact_telegram" binding:"required"`
	Medals            []int  `json:"medals" binding:"required"`
	Relative          string `json:"relative" binding:"required"`
}

type PersonDTO struct {
	ID                string     `json:"id"`
	Name              string     `json:"name"`
	Surname           string     `json:"surname"`
	Patronymic        string     `json:"patronymic"`
	DateBirth         int        `json:"date_birth"`
	DateDeath         int        `json:"date_death"`
	City              string     `json:"city"`
	History           string     `json:"history"`
	Rank              string     `json:"rank"`
	ContactEmail      string     `json:"contact_email"`
	ContactName       string     `json:"contact_name"`
	ContactSurname    string     `json:"contact_surname"`
	ContactPatronymic string     `json:"contact_patronymic"`
	ContactTelegram   string     `json:"contact_telegram"`
	StatusCheck       bool       `json:"status_check"`
	Medals            []MedalDTO `json:"medals"`
	Relative          string     `json:"relative"`
	DatePublished     time.Time  `json:"date_published"`
}

type PersonCountDTO struct {
	Count int `json:"count"`
}
