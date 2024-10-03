package db

import "time"

type Department struct {
	Id    *int32 `gorm:"primaryKey;autoIncrement"`
	Name  string
	Order int

	CreatedAt time.Time
	UpdatedAt int32
}

