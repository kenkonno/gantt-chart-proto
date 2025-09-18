package history

import (
	"github.com/kenkonno/gantt-chart-proto/backend/models/db"
	"github.com/kenkonno/gantt-chart-proto/backend/repository/connection"
	"github.com/kenkonno/gantt-chart-proto/backend/repository/interfaces"
	"gorm.io/gorm"
)

func NewHolidayRepository(historyId int32) interfaces.HolidayRepositoryIF {
	return &holidayRepository{
		historyId: historyId,
		con:       connection.GetCon(),
	}
}

type holidayRepository struct {
	historyId int32
	con       *gorm.DB
}

func (r *holidayRepository) FindAll() []db.Holiday {
	var holidays []db.Holiday
	result := r.con.Table("history_holidays").Where("history_id = ?", r.historyId).Order("id DESC").Find(&holidays)
	if result.Error != nil {
		panic(result.Error)
	}
	return holidays
}

func (r *holidayRepository) Find(id int32) db.Holiday {
	var holiday db.Holiday
	result := r.con.Table("history_holidays").Where("history_id = ? AND id = ?", r.historyId, id).First(&holiday)
	if result.Error != nil {
		panic(result.Error)
	}
	return holiday
}

func (r *holidayRepository) Upsert(m db.Holiday) {
	// History is read-only
}

func (r *holidayRepository) Delete(id int32) {
	// History is read-only
}
