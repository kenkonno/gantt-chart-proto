package repository
import (
	"github.com/kenkonno/gantt-chart-proto/backend/models/db"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)
// Auto generated start 
func NewOperationSettingRepository() operationSettingRepository {
	return operationSettingRepository{con}
}

type operationSettingRepository struct {
	con *gorm.DB
}
func (r *operationSettingRepository) FindAll() []db.OperationSetting {
	var operationSettings []db.OperationSetting

	result := r.con.Order("id DESC").Find(&operationSettings)
	if result.Error != nil {
		panic(result.Error)
	}
	return operationSettings
}

func (r *operationSettingRepository) Find(id int32) db.OperationSetting {
	var operationSetting db.OperationSetting

	result := r.con.First(&operationSetting, id)
	if result.Error != nil {
		panic(result.Error)
	}
	return operationSetting
}

func (r *operationSettingRepository) Upsert(m db.OperationSetting) {
	r.con.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		UpdateAll: true,
	}).Create(&m)
}

func (r *operationSettingRepository) Delete(id int32) {
	r.con.Where("id = ?", id).Delete(db.OperationSetting{})
}
// Auto generated end 
