package interfaces

import "github.com/kenkonno/gantt-chart-proto/backend/models/db"

type FacilityRepositoryIF interface {
	FindAll(facilityTypes []string, facilityStatus []string) []db.Facility
	Find(id int32) db.Facility
	Upsert(m db.Facility) db.Facility
	Delete(id int32)
}
