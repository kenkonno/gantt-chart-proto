package common
import (
	"github.com/kenkonno/gantt-chart-proto/backend/models/db"
	"github.com/kenkonno/gantt-chart-proto/backend/repository/connection"
	"github.com/kenkonno/gantt-chart-proto/backend/repository/interfaces"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)
// Auto generated start 
func NewMilestoneRepository() interfaces.MilestoneRepositoryIF {
	return &milestoneRepository{connection.GetCon()}
}

type milestoneRepository struct {
	con *gorm.DB
}
func (r *milestoneRepository) FindAll() []db.Milestone {
	var milestones []db.Milestone

	result := r.con.Order("id DESC").Find(&milestones)
	if result.Error != nil {
		panic(result.Error)
	}
	return milestones
}

func (r *milestoneRepository) FindByFacilityId(facilityId int32) []db.Milestone {
	var milestones []db.Milestone

	result := r.con.Order(`"order" DESC`).Where("facility_id = ?", facilityId).Find(&milestones)
	if result.Error != nil {
		panic(result.Error)
	}
	return milestones
}

func (r *milestoneRepository) Find(id int32) db.Milestone {
	var milestone db.Milestone

	result := r.con.First(&milestone, id)
	if result.Error != nil {
		panic(result.Error)
	}
	return milestone
}

func (r *milestoneRepository) Upsert(m db.Milestone) db.Milestone {
	r.con.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		UpdateAll: true,
	}).Create(&m)
	return m
}

func (r *milestoneRepository) Delete(id int32) {
	r.con.Where("id = ?", id).Delete(db.Milestone{})
}
// Auto generated end 
