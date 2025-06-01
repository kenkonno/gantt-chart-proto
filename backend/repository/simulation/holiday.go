package simulation

import (
	"github.com/kenkonno/gantt-chart-proto/backend/models/db"
	"github.com/kenkonno/gantt-chart-proto/backend/repository/connection"
	"github.com/kenkonno/gantt-chart-proto/backend/repository/interfaces"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Auto generated start
func NewSimulationHolidayRepository() interfaces.HolidayRepositoryIF {
	return &holidayRepository{
		con:   connection.GetCon(),
		table: "simulation_holidays",
	}
}

type holidayRepository struct {
	con   *gorm.DB
	table string
}

func (r *holidayRepository) FindAll() []db.Holiday {
	var holidays []db.Holiday

	result := r.con.Table(r.table).Order("id DESC").Find(&holidays)
	if result.Error != nil {
		panic(result.Error)
	}
	return holidays
}

func (r *holidayRepository) Find(id int32) db.Holiday {
	var holiday db.Holiday

	result := r.con.Table(r.table).First(&holiday, id)
	if result.Error != nil {
		panic(result.Error)
	}
	return holiday
}

func (r *holidayRepository) Upsert(m db.Holiday) {
	r.con.Table(r.table).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		UpdateAll: true,
	}).Create(&m)
}

func (r *holidayRepository) Delete(id int32) {
	r.con.Table(r.table).Where("id = ?", id).Delete(db.Holiday{})
}

// Auto generated end
