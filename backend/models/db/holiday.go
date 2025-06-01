package db

import "time"

type Holiday struct {
	Id   *int32 `gorm:"primaryKey;autoIncrement"`
	Name string
	Date time.Time

	CreatedAt time.Time
	UpdatedAt int32
}
