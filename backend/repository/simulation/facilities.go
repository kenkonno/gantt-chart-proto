package simulation

import (
	"github.com/kenkonno/gantt-chart-proto/backend/models/db"
	"github.com/kenkonno/gantt-chart-proto/backend/repository/connection"
	"github.com/kenkonno/gantt-chart-proto/backend/repository/interfaces"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Auto generated start
func NewSimulationFacilityRepository() interfaces.FacilityRepositoryIF {
	return &facilityRepository{
		con:   connection.GetCon(),
		table: "simulation_facilities",
	}
}

type facilityRepository struct {
	con *gorm.DB
	table string
}

func (r *facilityRepository) FindAll(facilityTypes []string, facilityStatus []string, orderColumn string) []db.Facility {
	var facilities []db.Facility

	builder := r.con.Table(r.table).Order(`"` + orderColumn + `" ASC`)
	if len(facilityTypes) > 0 {
		builder.Where("simulation_facilities.type IN ?", facilityTypes)
	}
	if len(facilityStatus) > 0 {
		builder.Where("simulation_facilities.status IN ?", facilityStatus)
	}

	result := builder.Find(&facilities)
	if result.Error != nil {
		panic(result.Error)
	}
	return facilities
}

func (r *facilityRepository) Find(id int32) db.Facility {
	var facility db.Facility

	result := r.con.Table(r.table).First(&facility, id)
	if result.Error != nil {
		panic(result.Error)
	}
	return facility
}

func (r *facilityRepository) Upsert(m db.Facility) db.Facility {
	r.con.Table(r.table).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		UpdateAll: true,
	}).Create(&m)
	return m
}

func (r *facilityRepository) Delete(id int32) {
	r.con.Table(r.table).Where("id = ?", id).Delete(db.Facility{})
}

// Auto generated end
