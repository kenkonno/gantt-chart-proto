package common

import (
	"time"

	"github.com/kenkonno/gantt-chart-proto/backend/models/db"
	"github.com/kenkonno/gantt-chart-proto/backend/repository/connection"
	"github.com/kenkonno/gantt-chart-proto/backend/repository/interfaces"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Auto generated start
func NewTicketDailyWeightRepository() interfaces.TicketDailyWeightRepositoryIF {
	return &ticketDailyWeightRepository{connection.GetCon()}
}

type ticketDailyWeightRepository struct {
	con *gorm.DB
}

func (r *ticketDailyWeightRepository) FindByFacilityId(facilityId int32) []db.TicketDailyWeight {
	var ticketDailyWeights []db.TicketDailyWeight

	result := r.con.
		Joins("INNER JOIN tickets ON ticket_daily_weights.ticket_id = tickets.id").
		Joins("INNER JOIN gantt_groups ON tickets.gantt_group_id = gantt_groups.id").
		Where("gantt_groups.facility_id = ?", facilityId).
		Order("ticket_daily_weights.ticket_id DESC").
		Find(&ticketDailyWeights)

	if result.Error != nil {
		panic(result.Error)
	}
	return ticketDailyWeights
}

func (r *ticketDailyWeightRepository) FindAll() []db.TicketDailyWeight {
	var ticketDailyWeights []db.TicketDailyWeight

	result := r.con.Order("ticket_id DESC").Find(&ticketDailyWeights)
	if result.Error != nil {
		panic(result.Error)
	}
	return ticketDailyWeights
}

func (r *ticketDailyWeightRepository) Find(ticketId int32, date time.Time) db.TicketDailyWeight {
	var ticketDailyWeight db.TicketDailyWeight

	result := r.con.Where("ticket_id = ? AND date = ?", ticketId, date).First(&ticketDailyWeight)
	if result.Error != nil {
		panic(result.Error)
	}
	return ticketDailyWeight
}
func (r *ticketDailyWeightRepository) FindByTicketId(ticketId int32) []db.TicketDailyWeight {
	var ticketDailyWeight []db.TicketDailyWeight

	result := r.con.Where("ticket_id = ?", ticketId).Find(&ticketDailyWeight)
	if result.Error != nil {
		panic(result.Error)
	}
	return ticketDailyWeight
}

func (r *ticketDailyWeightRepository) Upsert(m db.TicketDailyWeight) {
	r.con.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "ticket_id"}, {Name: "date"}},
		UpdateAll: true,
	}).Create(&m)
}

func (r *ticketDailyWeightRepository) Delete(ticketId int32, date time.Time) {
	r.con.Where("ticket_id = ? AND date = ?", ticketId, date).Delete(db.TicketDailyWeight{})
}

// Auto generated end
