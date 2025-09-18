package history

import (
	"github.com/kenkonno/gantt-chart-proto/backend/models/db"
	"github.com/kenkonno/gantt-chart-proto/backend/repository/connection"
	"github.com/kenkonno/gantt-chart-proto/backend/repository/interfaces"
	"gorm.io/gorm"
)

func NewProcessRepository(historyId int32) interfaces.ProcessRepositoryIF {
	return &processRepository{
		historyId: historyId,
		con:       connection.GetCon(),
	}
}

type processRepository struct {
	historyId int32
	con       *gorm.DB
}

func (r *processRepository) FindAll() []db.Process {
	var processes []db.Process
	result := r.con.Table("history_processes").Where("history_id = ?", r.historyId).Order(`"order" ASC`).Find(&processes)
	if result.Error != nil {
		panic(result.Error)
	}
	return processes
}

func (r *processRepository) Find(id int32) db.Process {
	var process db.Process
	result := r.con.Table("history_processes").Where("history_id = ? AND id = ?", r.historyId, id).First(&process)
	if result.Error != nil {
		panic(result.Error)
	}
	return process
}

func (r *processRepository) Upsert(m db.Process) {
	// History is read-only
}

func (r *processRepository) Delete(id int32) {
	// History is read-only
}
