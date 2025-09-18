package simulation

import (
	"github.com/kenkonno/gantt-chart-proto/backend/models/db"
	"github.com/kenkonno/gantt-chart-proto/backend/repository/connection"
	"github.com/kenkonno/gantt-chart-proto/backend/repository/interfaces"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Auto generated start
func NewSimulationFeatureOptionRepository() interfaces.FeatureOptionRepositoryIF {
	return &featureOptionRepository{
		con:   connection.GetCon(),
		table: "simulation_feature_options",
	}
}

type featureOptionRepository struct {
	con   *gorm.DB
	table string
}

func (r *featureOptionRepository) FindAll() []db.FeatureOption {
	var featureOptions []db.FeatureOption

	result := r.con.Table(r.table).Order("id DESC").Find(&featureOptions)
	if result.Error != nil {
		panic(result.Error)
	}
	return featureOptions
}

func (r *featureOptionRepository) Find(id int32) db.FeatureOption {
	var featureOption db.FeatureOption

	result := r.con.Table(r.table).First(&featureOption, id)
	if result.Error != nil {
		panic(result.Error)
	}
	return featureOption
}

func (r *featureOptionRepository) Upsert(m db.FeatureOption) {
	r.con.Table(r.table).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		UpdateAll: true,
	}).Create(&m)
}

func (r *featureOptionRepository) Delete(id int32) {
	r.con.Table(r.table).Where("id = ?", id).Delete(db.FeatureOption{})
}

// Auto generated end
