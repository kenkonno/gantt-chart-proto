package repository

import (
	"github.com/kenkonno/gantt-chart-proto/backend/models/db"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Auto generated start
func NewFacilityRepository() facilityRepository {
	return facilityRepository{con}
}

type facilityRepository struct {
	con *gorm.DB
}

func (r *facilityRepository) FindAll(facilityTypes []string, facilityStatus []string) []db.Facility {
	var facilities []db.Facility

	builder := r.con.Order(`"order" ASC`)
	if len(facilityTypes) > 0 {
		builder.Where("facilities.type IN ?", facilityTypes)
	}
	if len(facilityStatus) > 0 {
		builder.Where("facilities.status IN ?", facilityStatus)
	}

	result := builder.Find(&facilities)
	if result.Error != nil {
		panic(result.Error)
	}
	return facilities
}

func (r *facilityRepository) Find(id int32) db.Facility {
	var facility db.Facility

	result := r.con.First(&facility, id)
	if result.Error != nil {
		panic(result.Error)
	}
	return facility
}

func (r *facilityRepository) Upsert(m db.Facility) db.Facility {
	r.con.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		UpdateAll: true,
	}).Create(&m)
	return m
}

func (r *facilityRepository) Delete(id int32) {
	r.con.Where("id = ?", id).Delete(db.Facility{})
}

// Auto generated end
