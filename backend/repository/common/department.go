package common

import (
	"github.com/kenkonno/gantt-chart-proto/backend/models/db"
	"github.com/kenkonno/gantt-chart-proto/backend/repository/connection"
	"github.com/kenkonno/gantt-chart-proto/backend/repository/interfaces"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Auto generated start
func NewDepartmentRepository() interfaces.DepartmentRepositoryIF {
	return &departmentRepository{connection.GetCon()}
}

type departmentRepository struct {
	con *gorm.DB
}

func (r *departmentRepository) FindAll() []db.Department {
	var departments []db.Department

	result := r.con.Order(`"order" ASC`).Find(&departments)
	if result.Error != nil {
		panic(result.Error)
	}
	return departments
}

func (r *departmentRepository) Find(id int32) db.Department {
	var department db.Department

	result := r.con.First(&department, id)
	if result.Error != nil {
		panic(result.Error)
	}
	return department
}

func (r *departmentRepository) Upsert(m db.Department) {
	r.con.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns([]string{"name", "color", "order", "updated_at"}), // TODO: defaultを設定してるからかcolorが更新されなかった
	}).Debug().Create(&m)
}

func (r *departmentRepository) Delete(id int32) {
	r.con.Where("id = ?", id).Delete(db.Department{})
}

// Auto generated end
