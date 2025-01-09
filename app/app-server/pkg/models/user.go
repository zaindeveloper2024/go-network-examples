package models

type User struct {
	ID   string `json:"id"`
	Name string `json:"name" validate:"required,min=2"`
}
