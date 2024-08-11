package interfaces

import "github.com/kenkonno/gantt-chart-proto/backend/models/db"

type FacilitySharedLinkRepositoryIF interface {
	FindAll() []db.FacilitySharedLink
	Find(id int32) db.FacilitySharedLink
	FindByFacilityId(facilityId int32) *db.FacilitySharedLink
	Upsert(m db.FacilitySharedLink) db.FacilitySharedLink
	Delete(id int32)
	FindByUUID(uuid string) *db.FacilitySharedLink
}
