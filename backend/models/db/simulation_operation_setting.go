package db

import "time"

type SimulationOperationSetting struct {
	Id         *int32 `gorm:"primaryKey;autoIncrement"`
	FacilityId int32  `gorm:"uniqueIndex: operation_setting_u_index"`
	UnitId     int32  `gorm:"uniqueIndex: operation_setting_u_index"`
	ProcessId  int32  `gorm:"uniqueIndex: operation_setting_u_index"`
	WorkHour   int32

	CreatedAt time.Time
	UpdatedAt int32
}
