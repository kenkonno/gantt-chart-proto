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

func (r * simulationRepository) ResetSequence() {
	tableNames := []string{
		"departments",
		"facilities",
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
		simulationTableName := "simulation_" + tableName
		r.con.Exec(fmt.Sprintf("SELECT setval('%s', (SELECT MAX(id) FROM %s) + 1)",  tableName+"_id_seq", tableName))
		r.con.Exec(fmt.Sprintf("SELECT setval('%s', (SELECT MAX(id) FROM %s) + 1)",  simulationTableName+"_id_seq", simulationTableName))
	}
}

func (r * simulationRepository) SwitchTable() {
	tableNames := []string{
		"departments",
		"facilities",
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
		simulationTableName := "simulation_" + tableName
		r.con.Exec(fmt.Sprintf("ALTER TABLE %s RENAME TO %s",  tableName, tableName+"_temp"))
		r.con.Exec(fmt.Sprintf("ALTER TABLE %s RENAME TO %s",  simulationTableName, tableName))
		r.con.Exec(fmt.Sprintf("ALTER TABLE %s RENAME TO %s",  tableName+"_temp", simulationTableName))
	}
}