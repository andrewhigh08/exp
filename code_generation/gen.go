// Package entity — демо-модели для go generate (не package main, чтобы
// `go build ./...` в модуле code_generation не требовал func main).
//
// Директива go:generate ниже запускает генератор cmd/repogen через `go tool`:
// он подключен к модулю tool-директивой в go.mod (Go 1.24+), поэтому ставить
// repogen в PATH через go install не нужно. `go generate ./...` пересоздает
// gen_gen.go с UserRepository; этот файл коммитится вместе с исходниками.
package entity

//go:generate go tool repogen

//repogen:entity
type User struct {
	UserID       uint `gorm:"primaryKey"`
	Email        string
	PasswordHash string
}
