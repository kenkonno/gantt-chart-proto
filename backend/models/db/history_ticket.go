package db

import "time"

type HistoryTicket struct {
	Id              *int32 `gorm:"primaryKey"`
	HistoryId       int32  `gorm:"primaryKey"`
	GanttGroupId    int32
	ProcessId       *int32
	DepartmentId    *int32
	LimitDate       *time.Time
	Estimate        *int32
	NumberOfWorker  *int32
	DaysAfter       *int32
	StartDate       *time.Time
	EndDate         *time.Time
	ProgressPercent *int32
	Memo            string `gorm:"type:text"`
	Order           int32
	CreatedAt       time.Time
	UpdatedAt       int32
}

func (m *HistoryTicket) TableName() string {
	return "history_tickets"
}
