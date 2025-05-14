package interfaces

import "github.com/kenkonno/gantt-chart-proto/backend/models/db"

type FacilityWorkScheduleRepositoryIF interface {
	FindAll() []db.FacilityWorkSchedule
	Find(id int32) db.FacilityWorkSchedule
	Upsert(m db.FacilityWorkSchedule) db.FacilityWorkSchedule
	Delete(id int32)
}
