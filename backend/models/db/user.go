package db

import "time"

type User struct {
	ID       int32 `gorm:"primaryKey"`
	Password string
	Email    string

	CreatedAt time.Time
	UpdatedAt int
}
