// Package entity — демо-модели для go generate (не package main, чтобы
// `go build ./...` в модуле code_generation не требовал func main).
package entity

//go:generate repogen

//repogen:entity
type User struct {
	UserID       uint `gorm:"primary_key"`
	Email        string
	PasswordHash string
}
