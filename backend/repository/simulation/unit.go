package simulation

import (
	"github.com/kenkonno/gantt-chart-proto/backend/models/db"
	"github.com/kenkonno/gantt-chart-proto/backend/repository/connection"
	"github.com/kenkonno/gantt-chart-proto/backend/repository/interfaces"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Auto generated start
func NewSimulationUnitRepository() interfaces.UnitRepositoryIF {
	return &unitRepository{
		con:   connection.GetCon(),
		table: "simulation_units",
	}
}

type unitRepository struct {
	con *gorm.DB
	table string
}

func (r *unitRepository) FindAll() []db.Unit {
	var units []db.Unit

	result := r.con.Table(r.table).Order(`"order" ASC`).Find(&units)
	if result.Error != nil {
		panic(result.Error)
	}
	return units
}

func (r *unitRepository) Find(id int32) db.Unit {
	var unit db.Unit

	result := r.con.Table(r.table).First(&unit, id)
	if result.Error != nil {
		panic(result.Error)
	}
	return unit
}

func (r *unitRepository) Upsert(m db.Unit) db.Unit {
	r.con.Table(r.table).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		UpdateAll: true,
	}).Create(&m)
	return m
}

func (r *unitRepository) Delete(id int32) {
	r.con.Table(r.table).Where("id = ?", id).Delete(db.Unit{})
}

// Auto generated end
func (r *unitRepository) FindByFacilityId(facilityId int32) []db.Unit {
	var units []db.Unit

	result := r.con.Table(r.table).Where("facility_id = ?", facilityId).Order("simulation_units.order ASC").Find(&units)
	if result.Error != nil {
		panic(result.Error)
	}
	return units
}
