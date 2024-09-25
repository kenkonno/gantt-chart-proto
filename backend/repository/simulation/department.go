package simulation

import (
	"github.com/kenkonno/gantt-chart-proto/backend/models/db"
	"github.com/kenkonno/gantt-chart-proto/backend/repository/connection"
	"github.com/kenkonno/gantt-chart-proto/backend/repository/interfaces"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Auto generated start
func NewSimulationDepartmentRepository() interfaces.DepartmentRepositoryIF {
	return &departmentRepository{
		con:   connection.GetCon(),
		table: "simulation_departments",
	}
}

type departmentRepository struct {
	con   *gorm.DB
	table string
}

func (r *departmentRepository) FindAll() []db.Department {
	var departments []db.Department

	result := r.con.Table(r.table).Order(`"order" ASC`).Find(&departments)
	if result.Error != nil {
		panic(result.Error)
	}
	return departments
}

func (r *departmentRepository) Find(id int32) db.Department {
	var department db.Department

	result := r.con.Table(r.table).First(&department, id)
	if result.Error != nil {
		panic(result.Error)
	}
	return department
}

func (r *departmentRepository) Upsert(m db.Department) {
	r.con.Table(r.table).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		UpdateAll: true,
	}).Create(&m)
}

func (r *departmentRepository) Delete(id int32) {
	r.con.Table(r.table).Where("id = ?", id).Delete(db.Department{})
}

// Auto generated end
