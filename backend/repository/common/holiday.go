package common

import (
	"github.com/kenkonno/gantt-chart-proto/backend/models/db"
	"github.com/kenkonno/gantt-chart-proto/backend/repository/connection"
	"github.com/kenkonno/gantt-chart-proto/backend/repository/interfaces"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Auto generated start
func NewHolidayRepository() interfaces.HolidayRepositoryIF {
	return &holidayRepository{connection.GetCon()}
}

type holidayRepository struct {
	con *gorm.DB
}

func (r *holidayRepository) FindAll() []db.Holiday {
	var holidays []db.Holiday

	result := r.con.Order("date DESC").Find(&holidays)
	if result.Error != nil {
		panic(result.Error)
	}
	return holidays
}

func (r *holidayRepository) Find(id int32) db.Holiday {
	var holiday db.Holiday

	result := r.con.First(&holiday, id)
	if result.Error != nil {
		panic(result.Error)
	}
	return holiday
}

func (r *holidayRepository) Upsert(m db.Holiday) {
	r.con.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		UpdateAll: true,
	}).Create(&m)
}

func (r *holidayRepository) Delete(id int32) {
	r.con.Where("id = ?", id).Delete(db.Holiday{})
}

// Auto generated end
