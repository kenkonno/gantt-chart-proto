package history

import (
	"github.com/kenkonno/gantt-chart-proto/backend/models/db"
	"github.com/kenkonno/gantt-chart-proto/backend/repository/connection"
	"github.com/kenkonno/gantt-chart-proto/backend/repository/interfaces"
	"gorm.io/gorm"
)

func NewDepartmentRepository(historyId int32) interfaces.DepartmentRepositoryIF {
	return &departmentRepository{
		historyId: historyId,
		con:       connection.GetCon(),
	}
}

type departmentRepository struct {
	historyId int32
	con       *gorm.DB
}

func (r *departmentRepository) FindAll() []db.Department {
	var departments []db.Department
	result := r.con.Table("history_departments").Where("history_id = ?", r.historyId).Order(`"order" ASC`).Find(&departments)
	if result.Error != nil {
		panic(result.Error)
	}
	return departments
}

func (r *departmentRepository) Find(id int32) db.Department {
	var department db.Department
	result := r.con.Table("history_departments").Where("history_id = ? AND id = ?", r.historyId, id).First(&department)
	if result.Error != nil {
		panic(result.Error)
	}
	return department
}

func (r *departmentRepository) Upsert(m db.Department) {
	// History is read-only
}

func (r *departmentRepository) Delete(id int32) {
	// History is read-only
}
