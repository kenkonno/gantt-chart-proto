package interfaces

import (
	"github.com/kenkonno/gantt-chart-proto/backend/models/db"
	"time"
)

type HolidayRepositoryIF interface {
	FindAll() []db.Holiday
	Find(id int32) db.Holiday
	Upsert(m db.Holiday)
	Delete(id int32)
	FindByFacilityId(facilityId int32) []db.Holiday
	InsertByFacilityId(facilityId int32, from *time.Time, to *time.Time) []db.Holiday
}
