package simulation

import (
	"github.com/kenkonno/gantt-chart-proto/backend/models/db"
	"github.com/kenkonno/gantt-chart-proto/backend/repository/connection"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Auto generated start
func NewSimulationLock() *simulationLockRepository {
	return &simulationLockRepository{
		con: connection.GetCon(),
	}
}

type simulationLockRepository struct {
	con *gorm.DB
}

func (r *simulationLockRepository) FindAll() []db.SimulationLock {
	var simulationLocks []db.SimulationLock

	result := r.con.Order(`"order" ASC`).Find(&simulationLocks)
	if result.Error != nil {
		panic(result.Error)
	}
	return simulationLocks
}

func (r *simulationLockRepository) Find(simulationName string) db.SimulationLock {
	var simulationLock db.SimulationLock

	result := r.con.First(&simulationLock, simulationName)
	if result.Error != nil {
		panic(result.Error)
	}
	return simulationLock
}

func (r *simulationLockRepository) Upsert(m db.SimulationLock) {
	r.con.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		UpdateAll: true,
	}).Create(&m)
}

func (r *simulationLockRepository) Delete(simulationName string) {
	r.con.Where("simulation_name = ?", simulationName).Delete(db.SimulationLock{})
}

// Auto generated end
