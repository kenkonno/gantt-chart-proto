package db

import "time"

type Process struct {
	Id    *int32 `gorm:"primaryKey;autoIncrement"`
	Name  string
	Order int

	CreatedAt time.Time
	UpdatedAt int
}
