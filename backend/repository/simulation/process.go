package simulation

import (
	"github.com/kenkonno/gantt-chart-proto/backend/models/db"
	"github.com/kenkonno/gantt-chart-proto/backend/repository/connection"
	"github.com/kenkonno/gantt-chart-proto/backend/repository/interfaces"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Auto generated start
func NewSimulationProcessRepository() interfaces.ProcessRepositoryIF {
	return &processRepository{
		con:   connection.GetCon(),
		table: "simulation_processes",
	}
}

type processRepository struct {
	con *gorm.DB
	table string
}

func (r *processRepository) FindAll() []db.Process {
	var processes []db.Process

	result := r.con.Table(r.table).Order(`"order" ASC`).Find(&processes)
	if result.Error != nil {
		panic(result.Error)
	}
	return processes
}

func (r *processRepository) Find(id int32) db.Process {
	var process db.Process

	result := r.con.Table(r.table).First(&process, id)
	if result.Error != nil {
		panic(result.Error)
	}
	return process
}

func (r *processRepository) Upsert(m db.Process) {
	r.con.Table(r.table).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		UpdateAll: true,
	}).Create(&m)
}

func (r *processRepository) Delete(id int32) {
	r.con.Table(r.table).Where("id = ?", id).Delete(db.Process{})
}

// Auto generated end
