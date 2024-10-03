package db

import (
	"time"
)

type SimulationFacilitySharedLink struct {
	Id         *int32 `gorm:"primaryKey;autoIncrement"`
	FacilityId int32
	Uuid       string
	CreatedAt  time.Time
	UpdatedAt  int32
}
