package db

import "time"

type OperationSetting struct {
	Id         *int32 `gorm:"primaryKey;autoIncrement"`
	FacilityId int32
	UnitId     int32
	ProcessId  int32
	WorkHour   int32

	CreatedAt time.Time
	UpdatedAt int
}
