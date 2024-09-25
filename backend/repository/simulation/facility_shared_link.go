package simulation

import (
	"github.com/kenkonno/gantt-chart-proto/backend/models/db"
	"github.com/kenkonno/gantt-chart-proto/backend/repository/connection"
	"github.com/kenkonno/gantt-chart-proto/backend/repository/interfaces"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Auto generated start
func NewSimulationFacilitySharedLinkRepository() interfaces.FacilitySharedLinkRepositoryIF {
	return &facilitySharedLinkRepository{
		con:   connection.GetCon(),
		table: "simulation_departments",
	}
}

type facilitySharedLinkRepository struct {
	con *gorm.DB
	table string
}

func (r *facilitySharedLinkRepository) FindAll() []db.FacilitySharedLink {
	var facilitySharedLinks []db.FacilitySharedLink

	result := r.con.Table(r.table).Order("id DESC").Find(&facilitySharedLinks)
	if result.Error != nil {
		panic(result.Error)
	}
	return facilitySharedLinks
}

func (r *facilitySharedLinkRepository) Find(id int32) db.FacilitySharedLink {
	var facilitySharedLink db.FacilitySharedLink

	result := r.con.Table(r.table).First(&facilitySharedLink, id)
	if result.Error != nil {
		panic(result.Error)
	}
	return facilitySharedLink
}
func (r *facilitySharedLinkRepository) FindByFacilityId(facilityId int32) *db.FacilitySharedLink {
	var facilitySharedLink *db.FacilitySharedLink

	result := r.con.Table(r.table).Where("facility_id", facilityId).First(&facilitySharedLink)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			// nothing todo...
		} else {
			panic(result.Error)
		}
	}
	return facilitySharedLink
}

func (r *facilitySharedLinkRepository) Upsert(m db.FacilitySharedLink) db.FacilitySharedLink {
	r.con.Table(r.table).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		UpdateAll: true,
	}).Create(&m)
	return m
}

func (r *facilitySharedLinkRepository) Delete(id int32) {
	r.con.Table(r.table).Where("id = ?", id).Delete(db.FacilitySharedLink{})
}

func (r *facilitySharedLinkRepository) FindByUUID(uuid string) *db.FacilitySharedLink {
	var facilitySharedLink *db.FacilitySharedLink

	result := r.con.Table(r.table).Where("uuid", uuid).First(&facilitySharedLink)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			// nothing todo...
		} else {
			panic(result.Error)
		}
	}
	return facilitySharedLink
}

// Auto generated end
