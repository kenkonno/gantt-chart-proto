package db

import "time"

type Unit struct {
	Id         *int32 `gorm:"primaryKey;autoIncrement"`
	Name       string
	FacilityId int32

	Order     int
	CreatedAt time.Time
	UpdatedAt int
}
