package repository

import (
	"github.com/kenkonno/gantt-chart-proto/backend/models/db"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Auto generated start
func NewUnitRepository() unitRepository {
	return unitRepository{con}
}

type unitRepository struct {
	con *gorm.DB
}

func (r *unitRepository) FindAll() []db.Unit {
	var units []db.Unit

	result := r.con.Order("id DESC").Find(&units)
	if result.Error != nil {
		panic(result.Error)
	}
	return units
}

func (r *unitRepository) Find(id int32) db.Unit {
	var unit db.Unit

	result := r.con.First(&unit, id)
	if result.Error != nil {
		panic(result.Error)
	}
	return unit
}

func (r *unitRepository) Upsert(m db.Unit) {
	r.con.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		UpdateAll: true,
	}).Create(&m)
}

func (r *unitRepository) Delete(id int32) {
	r.con.Where("id = ?", id).Delete(db.Unit{})
}

// Auto generated end
func (r *unitRepository) FindByFacilityId(facilityId int32) []db.Unit {
	var units []db.Unit

	result := r.con.Where("facility_id = ?", facilityId).Order("id DESC").Find(&units)
	if result.Error != nil {
		panic(result.Error)
	}
	return units
}
