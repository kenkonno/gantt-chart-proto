package db

import "time"

type SimulationGanttGroup struct {
	Id         *int32 `gorm:"primaryKey;autoIncrement"`
	FacilityId int32
	UnitId     int32

	CreatedAt time.Time
	UpdatedAt int32
}
