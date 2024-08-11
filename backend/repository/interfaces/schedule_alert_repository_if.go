package interfaces

import "github.com/kenkonno/gantt-chart-proto/backend/models/db"

type ScheduleAlertRepositoryIF interface {
	FindAll() []db.ScheduleAlert
	Find(id int32) db.ScheduleAlert
}
