package main

//go:generate repogen

//repogen:entity
type User struct {
	UserID       uint `gorm:"primary_key"`
	Email        string
	PasswordHash string
}
