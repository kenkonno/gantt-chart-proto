package db

import "time"

type Department struct {
	Id    *int32 `gorm:"primaryKey;autoIncrement"`
	Name  string
	Color string `gorm:"default:'rgb(66, 165, 246)'"`
	Order int

	CreatedAt time.Time
	UpdatedAt int32
}

