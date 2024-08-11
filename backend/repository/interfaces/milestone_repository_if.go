package interfaces

import "github.com/kenkonno/gantt-chart-proto/backend/models/db"

type MilestoneRepositoryIF interface {
	FindAll() []db.Milestone
	FindByFacilityId(facilityId int32) []db.Milestone
	Find(id int32) db.Milestone
	Upsert(m db.Milestone) db.Milestone
	Delete(id int32)
}
