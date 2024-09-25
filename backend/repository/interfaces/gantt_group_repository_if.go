package interfaces

import "github.com/kenkonno/gantt-chart-proto/backend/models/db"

type GanttGroupRepositoryIF interface {
	FindAll() []db.GanttGroup
	Find(id int32) db.GanttGroup
	Upsert(m db.GanttGroup) db.GanttGroup
	Delete(id int32)
	FindByFacilityId(facilityId int32) []db.GanttGroup
	DeleteByUnitId(unitId int32)
}
