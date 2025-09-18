package history

import (
	"github.com/kenkonno/gantt-chart-proto/backend/models/db"
	"github.com/kenkonno/gantt-chart-proto/backend/repository/connection"
	"github.com/kenkonno/gantt-chart-proto/backend/repository/interfaces"
	"gorm.io/gorm"
)

func NewUnitRepository(historyId int32) interfaces.UnitRepositoryIF {
	return &unitRepository{
		historyId: historyId,
		con:       connection.GetCon(),
	}
}

type unitRepository struct {
	historyId int32
	con       *gorm.DB
}

func (r *unitRepository) FindAll() []db.Unit {
	var units []db.Unit
	result := r.con.Table("history_units").Where("history_id = ?", r.historyId).Order(`"order" ASC`).Find(&units)
	if result.Error != nil {
		panic(result.Error)
	}
	return units
}

func (r *unitRepository) Find(id int32) db.Unit {
	var unit db.Unit
	result := r.con.Table("history_units").Where("history_id = ? AND id = ?", r.historyId, id).First(&unit)
	if result.Error != nil {
		panic(result.Error)
	}
	return unit
}

func (r *unitRepository) Upsert(m db.Unit) db.Unit {
	// History is read-only
	return m
}

func (r *unitRepository) Delete(id int32) {
	// History is read-only
}

func (r *unitRepository) FindByFacilityId(facilityId int32) []db.Unit {
	var units []db.Unit
	result := r.con.Table("history_units").Where("history_id = ? AND facility_id = ?", r.historyId, facilityId).Order("history_units.order ASC").Find(&units)
	if result.Error != nil {
		panic(result.Error)
	}
	return units
}
