package db

import "time"

type Process struct {
	Id   *int32 `gorm:"primaryKey;autoIncrement"`
	Name string

	CreatedAt time.Time
	UpdatedAt int
}
