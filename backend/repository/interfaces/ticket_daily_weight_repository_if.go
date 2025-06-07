package interfaces

import (
	"github.com/kenkonno/gantt-chart-proto/backend/models/db"
)

type TicketDailyWeightRepositoryIF interface {
	FindAll() []db.TicketDailyWeight
	Find(id int32) db.TicketDailyWeight
	Upsert(m db.TicketDailyWeight)
	Delete(id int32)
}
