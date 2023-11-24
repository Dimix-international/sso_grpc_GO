package models

type App struct {
	ID     int
	Name   string
	Secret string //для подписи токенов
}