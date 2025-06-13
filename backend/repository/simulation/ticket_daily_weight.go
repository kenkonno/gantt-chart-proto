package simulation

import (
	"github.com/kenkonno/gantt-chart-proto/backend/models/db"
	"github.com/kenkonno/gantt-chart-proto/backend/repository/connection"
	"github.com/kenkonno/gantt-chart-proto/backend/repository/interfaces"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

// Auto generated start
func NewSimulationTicketDailyWeightRepository() interfaces.TicketDailyWeightRepositoryIF {
	return &ticketDailyWeightRepository{
		con:   connection.GetCon(),
		table: "simulation_ticket_daily_weights",
	}
}

type ticketDailyWeightRepository struct {
	con   *gorm.DB
	table string
}

func (r *ticketDailyWeightRepository) FindAll() []db.TicketDailyWeight {
	var ticketDailyWeights []db.TicketDailyWeight

	result := r.con.Table(r.table).Order("ticket_id DESC").Find(&ticketDailyWeights)
	if result.Error != nil {
		panic(result.Error)
	}
	return ticketDailyWeights
}

func (r *ticketDailyWeightRepository) Find(ticketId int32, date time.Time) db.TicketDailyWeight {
	var ticketDailyWeight db.TicketDailyWeight

	result := r.con.Table(r.table).Where("ticket_id = ? AND date = ?", ticketId, date).First(&ticketDailyWeight)
	if result.Error != nil {
		panic(result.Error)
	}
	return ticketDailyWeight
}

func (r *ticketDailyWeightRepository) FindByTicketId(ticketId int32) []db.TicketDailyWeight {
	var ticketDailyWeight []db.TicketDailyWeight

	result := r.con.Table(r.table).Where("ticket_id = ?", ticketId).Find(&ticketDailyWeight)
	if result.Error != nil {
		panic(result.Error)
	}
	return ticketDailyWeight
}

func (r *ticketDailyWeightRepository) Upsert(m db.TicketDailyWeight) {
	r.con.Table(r.table).Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "ticket_id"}, {Name: "date"}},
		UpdateAll: true,
	}).Create(&m)
}

func (r *ticketDailyWeightRepository) Delete(ticketId int32, date time.Time) {
	r.con.Table(r.table).Where("ticket_id = ? AND date = ?", ticketId, date).Delete(db.TicketDailyWeight{})
}

// Auto generated end
