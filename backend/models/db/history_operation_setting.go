package db

import "time"

type HistoryOperationSetting struct {
	Id         *int32 `gorm:"primaryKey"`
	HistoryId  int32  `gorm:"primaryKey;uniqueIndex:history_operation_setting_u_index"`
	FacilityId int32  `gorm:"uniqueIndex:history_operation_setting_u_index"`
	UnitId     int32  `gorm:"uniqueIndex:history_operation_setting_u_index"`
	ProcessId  int32  `gorm:"uniqueIndex:history_operation_setting_u_index"`
	WorkHour   int32
	CreatedAt  time.Time
	UpdatedAt  int32
}

func (m *HistoryOperationSetting) TableName() string {
	return "history_operation_settings"
}
