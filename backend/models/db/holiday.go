package db

import "time"

type Holiday struct {
	Id         *int32 `gorm:"primaryKey;autoIncrement"`
	FacilityId int32
	Name       string
	Date       time.Time

	CreatedAt time.Time
	UpdatedAt int32
}
