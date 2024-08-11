package interfaces

import "github.com/kenkonno/gantt-chart-proto/backend/models/db"

type TicketRepositoryIF interface {
	FindAll() []db.Ticket
	FindByFacilityType(facilityTypes []string, facilityStatus []string) []db.Ticket
	Find(id int32) db.Ticket
	Upsert(m db.Ticket) (db.Ticket, error)
	UpdateMemo(m db.Ticket) (db.Ticket, error)
	Delete(id int32)
	FindByGanttGroupIds(ganttGroupIds []int32) []db.Ticket
}
