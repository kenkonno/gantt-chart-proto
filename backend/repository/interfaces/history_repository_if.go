package interfaces

import "github.com/kenkonno/gantt-chart-proto/backend/models/db"

type HistoryRepositoryIF interface {
	CreateSnapshot(facilityId int32, name string) (int32, error)
	FindByFacilityId(facilityId int32) []db.HistoryFacility
	UpdateName(id int32, name string) error
	Delete(id int32) error
}
