package simulation

import (
	"fmt"
	"github.com/kenkonno/gantt-chart-proto/backend/repository/connection"
	"gorm.io/gorm"
)

// シミュレーション専用のRepository,初期化や例外的な操作を実施する。
func NewSimulationRepository() *simulationRepository {
	return &simulationRepository{
		con: connection.GetCon(),
	}
}

type simulationRepository struct {
	con *gorm.DB
}

func (r *simulationRepository) InitAllData() {
	// tables
	tableNames := []string{
		"departments",
		"facilities",
		"facility_shared_links",
		"gantt_groups",
		"holidays",
		"milestones",
		"operation_settings",
		"processes",
		"tickets",
		"ticket_users",
		"units",
		"users",
	}

	// Truncate, Copy
	for _, tableName := range tableNames {
		r.con.Exec(fmt.Sprintf("TRUNCATE TABLE %s", "simulation_" + tableName))
		r.con.Exec(fmt.Sprintf("INSERT INTO %s SELECT * FROM %s", "simulation_" + tableName, tableName))
	}
}
