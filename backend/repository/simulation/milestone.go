package simulation
import (
	"github.com/kenkonno/gantt-chart-proto/backend/models/db"
	"github.com/kenkonno/gantt-chart-proto/backend/repository/connection"
	"github.com/kenkonno/gantt-chart-proto/backend/repository/interfaces"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)
// Auto generated start 
func NewSimulationMilestoneRepository() interfaces.MilestoneRepositoryIF {
	return &milestoneRepository{
		con:   connection.GetCon(),
		table: "simulation_milestones",
	}
}

type milestoneRepository struct {
	con *gorm.DB
	table string
}
func (r *milestoneRepository) FindAll() []db.Milestone {
	var milestones []db.Milestone

	result := r.con.Table(r.table).Order("id DESC").Find(&milestones)
	if result.Error != nil {
		panic(result.Error)
	}
	return milestones
}

func (r *milestoneRepository) FindByFacilityId(facilityIds []int32) []db.Milestone {
	var milestones []db.Milestone

	result := r.con.Table(r.table).Order(`"order" DESC`).Where("facility_id IN ?", facilityIds).Find(&milestones)
	if result.Error != nil {
		panic(result.Error)
	}
	return milestones
}

func (r *milestoneRepository) Find(id int32) db.Milestone {
	var milestone db.Milestone

	result := r.con.Table(r.table).First(&milestone, id)
	if result.Error != nil {
		panic(result.Error)
	}
	return milestone
}

func (r *milestoneRepository) Upsert(m db.Milestone) db.Milestone {
	r.con.Table(r.table).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		UpdateAll: true,
	}).Create(&m)
	return m
}

func (r *milestoneRepository) Delete(id int32) {
	r.con.Table(r.table).Where("id = ?", id).Delete(db.Milestone{})
}
// Auto generated end 
