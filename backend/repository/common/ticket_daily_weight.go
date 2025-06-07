package common

import (
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

func (r *ticketDailyWeightRepository) FindAll() []db.TicketDailyWeight {
	var ticketDailyWeights []db.TicketDailyWeight

	result := r.con.Order("id DESC").Find(&ticketDailyWeights)
	if result.Error != nil {
		panic(result.Error)
	}
	return ticketDailyWeights
}

func (r *ticketDailyWeightRepository) Find(id int32) db.TicketDailyWeight {
	var ticketDailyWeight db.TicketDailyWeight

	result := r.con.First(&ticketDailyWeight, id)
	if result.Error != nil {
		panic(result.Error)
	}
	return ticketDailyWeight
}

func (r *ticketDailyWeightRepository) Upsert(m db.TicketDailyWeight) {
	r.con.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		UpdateAll: true,
	}).Create(&m)
}

func (r *ticketDailyWeightRepository) Delete(id int32) {
	r.con.Where("id = ?", id).Delete(db.TicketDailyWeight{})
}

// Auto generated end
