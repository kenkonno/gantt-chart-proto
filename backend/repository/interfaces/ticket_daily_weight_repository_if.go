package interfaces

import (
	"time"

	"github.com/kenkonno/gantt-chart-proto/backend/models/db"
)

type TicketDailyWeightRepositoryIF interface {
	FindAll() []db.TicketDailyWeight
	FindByFacilityId(facilityId int32) []db.TicketDailyWeight
	Find(ticketId int32, date time.Time) db.TicketDailyWeight
	FindByTicketId(ticketId int32) []db.TicketDailyWeight
	Upsert(m db.TicketDailyWeight)
	Delete(ticketId int32, date time.Time)
}
