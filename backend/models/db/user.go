package db

import "time"

type User struct {
	Id               *int32 `gorm:"primaryKey;autoIncrement"`
	DepartmentId     int32
	LimitOfOperation float32
	Name             string
	Password         string
	Email            string

	CreatedAt time.Time
	UpdatedAt int64
}
