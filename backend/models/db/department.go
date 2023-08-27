package db

import "time"

type Department struct {
	Id   *int32 `gorm:"primaryKey;autoIncrement"`
	Name string

	CreatedAt time.Time
	UpdatedAt int
}
