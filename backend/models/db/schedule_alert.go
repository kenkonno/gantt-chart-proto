package db

import "time"

type ScheduleAlert struct {
	FacilityId         int32
	FacilityName       string
	UnitId             int32
	UnitName           string
	ProcessId          int32
	ProcessName        string
	StartDate          time.Time
	EndDate            time.Time
	ProgressPercent    float32
	ActualProgressDate time.Time // 進捗から計算した完了日
	DelayDays          float32   // 終わっていない工数h
	TicketId           int32
}
