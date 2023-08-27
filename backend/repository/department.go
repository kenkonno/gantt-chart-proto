package repository
import (
	"github.com/kenkonno/gantt-chart-proto/backend/models/db"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)
// Auto generated start 
func NewDepartmentRepository() departmentRepository {
	return departmentRepository{con}
}

type departmentRepository struct {
	con *gorm.DB
}
func (r *departmentRepository) FindAll() []db.Department {
	var departments []db.Department

	result := r.con.Order("id DESC").Find(&departments)
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
		UpdateAll: true,
	}).Create(&m)
}

func (r *departmentRepository) Delete(id int32) {
	r.con.Where("id = ?", id).Delete(db.Department{})
}
// Auto generated end 
