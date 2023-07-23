package db

import "time"

type User struct {
	Id       int32 `gorm:"primaryKey"`
	Password string
	Email    string

	CreatedAt time.Time
	UpdatedAt int
}
