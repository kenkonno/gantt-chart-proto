package history

import (
	"github.com/kenkonno/gantt-chart-proto/backend/models/db"
	"github.com/kenkonno/gantt-chart-proto/backend/repository/connection"
	"github.com/kenkonno/gantt-chart-proto/backend/repository/interfaces"
	"gorm.io/gorm"
)

func NewMilestoneRepository(historyId int32) interfaces.MilestoneRepositoryIF {
	return &milestoneRepository{
		historyId: historyId,
		con:       connection.GetCon(),
	}
}

type milestoneRepository struct {
	historyId int32
	con       *gorm.DB
}

func (r *milestoneRepository) FindAll() []db.Milestone {
	var milestones []db.Milestone
	result := r.con.Table("history_milestones").Where("history_id = ?", r.historyId).Order("id DESC").Find(&milestones)
	if result.Error != nil {
		panic(result.Error)
	}
	return milestones
}

func (r *milestoneRepository) FindByFacilityId(facilityIds []int32) []db.Milestone {
	var milestones []db.Milestone
	result := r.con.Table("history_milestones").Where("history_id = ? AND facility_id IN ?", r.historyId, facilityIds).Order(`"order" DESC`).Find(&milestones)
	if result.Error != nil {
		panic(result.Error)
	}
	return milestones
}

func (r *milestoneRepository) Find(id int32) db.Milestone {
	var milestone db.Milestone
	result := r.con.Table("history_milestones").Where("history_id = ? AND id = ?", r.historyId, id).First(&milestone)
	if result.Error != nil {
		panic(result.Error)
	}
	return milestone
}

func (r *milestoneRepository) Upsert(m db.Milestone) db.Milestone {
	// History is read-only
	return m
}

func (r *milestoneRepository) Delete(id int32) {
	// History is read-only
}
