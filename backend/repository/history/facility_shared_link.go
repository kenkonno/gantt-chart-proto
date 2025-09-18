package history

import (
	"github.com/kenkonno/gantt-chart-proto/backend/models/db"
	"github.com/kenkonno/gantt-chart-proto/backend/repository/connection"
	"github.com/kenkonno/gantt-chart-proto/backend/repository/interfaces"
	"gorm.io/gorm"
)

func NewFacilitySharedLinkRepository(historyId int32) interfaces.FacilitySharedLinkRepositoryIF {
	return &facilitySharedLinkRepository{
		historyId: historyId,
		con:       connection.GetCon(),
	}
}

type facilitySharedLinkRepository struct {
	historyId int32
	con       *gorm.DB
}

func (r *facilitySharedLinkRepository) FindAll() []db.FacilitySharedLink {
	var facilitySharedLinks []db.FacilitySharedLink
	result := r.con.Table("history_facility_shared_links").Where("history_id = ?", r.historyId).Find(&facilitySharedLinks)
	if result.Error != nil {
		panic(result.Error)
	}
	return facilitySharedLinks
}

func (r *facilitySharedLinkRepository) Find(id int32) db.FacilitySharedLink {
	var facilitySharedLink db.FacilitySharedLink
	result := r.con.Table("history_facility_shared_links").Where("history_id = ? AND id = ?", r.historyId, id).First(&facilitySharedLink)
	if result.Error != nil {
		panic(result.Error)
	}
	return facilitySharedLink
}

func (r *facilitySharedLinkRepository) FindByFacilityId(facilityId int32) *db.FacilitySharedLink {
	var facilitySharedLink db.FacilitySharedLink
	result := r.con.Table("history_facility_shared_links").Where("history_id = ? AND facility_id = ?", r.historyId, facilityId).First(&facilitySharedLink)
	if result.Error != nil {
		return nil
	}
	return &facilitySharedLink
}

func (r *facilitySharedLinkRepository) Upsert(m db.FacilitySharedLink) db.FacilitySharedLink {
	// History is read-only
	return m
}

func (r *facilitySharedLinkRepository) Delete(id int32) {
	// History is read-only
}

func (r *facilitySharedLinkRepository) FindByUUID(uuid string) *db.FacilitySharedLink {
	var facilitySharedLink db.FacilitySharedLink
	result := r.con.Table("history_facility_shared_links").Where("history_id = ? AND uuid = ?", r.historyId, uuid).First(&facilitySharedLink)
	if result.Error != nil {
		return nil
	}
	return &facilitySharedLink
}
