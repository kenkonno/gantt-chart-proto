package interfaces

import "github.com/kenkonno/gantt-chart-proto/backend/models/db"

type UnitRepositoryIF interface {
	FindAll() []db.Unit
	Find(id int32) db.Unit
	Upsert(m db.Unit) db.Unit
	Delete(id int32)
	FindByFacilityId(facilityId int32) []db.Unit
}
