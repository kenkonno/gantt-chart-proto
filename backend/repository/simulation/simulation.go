package simulation

import (
	"fmt"
	"github.com/kenkonno/gantt-chart-proto/backend/repository/connection"
	"github.com/samber/lo"
	"gorm.io/gorm"
	"strings"
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

	columns := [][]string{
		{"id", "name", "order", "created_at", "updated_at"},
		{"id", "name", "term_from", "term_to", "order", "created_at", "updated_at", "status", "type", "shipment_due_date"},
		{"id", "facility_id", "unit_id", "created_at", "updated_at"},
		{"id", "facility_id", "name", "date", "created_at", "updated_at"},
		{"id", "facility_id", "date", "description", "order", "created_at", "updated_at"},
		{"id", "facility_id", "unit_id", "process_id", "work_hour", "created_at", "updated_at"},
		{"id", "name", "order", "created_at", "updated_at", "color"},
		{"id", "gantt_group_id", "process_id", "department_id", "limit_date", "estimate", "number_of_worker", "days_after", "start_date", "end_date", "progress_percent", "order", "created_at", "updated_at", "memo"},
		{"id", "ticket_id", "user_id", "order", "created_at", "updated_at"},
		{"id", "name", "facility_id", "order", "created_at", "updated_at"},
		{"id", "department_id", "limit_of_operation", "password", "email", "created_at", "updated_at", "role", "last_name", "first_name", "password_reset"},
	}

	// Truncate, Copy
	for index, tableName := range tableNames {
		column := strings.Join(lo.Map(columns[index], func(item string, index int) string {
			return `"` + item + `"`
		}), ",")
		r.con.Exec(fmt.Sprintf("TRUNCATE TABLE %s", "simulation_"+tableName))
		r.con.Exec(fmt.Sprintf("INSERT INTO %s (%s) SELECT %s FROM %s", "simulation_"+tableName, column, column, tableName))
	}
}

func (r *simulationRepository) ResetSequence() {
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
		r.con.Exec(fmt.Sprintf("SELECT setval('%s', (SELECT MAX(id) FROM %s) + 1)", tableName+"_id_seq", tableName))
		r.con.Exec(fmt.Sprintf("SELECT setval('%s', (SELECT MAX(id) FROM %s) + 1)", simulationTableName+"_id_seq", simulationTableName))
	}
}

func (r *simulationRepository) SwitchTable() {
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
		r.con.Exec(fmt.Sprintf("ALTER TABLE %s RENAME TO %s", tableName, tableName+"_temp"))
		r.con.Exec(fmt.Sprintf("ALTER TABLE %s RENAME TO %s", simulationTableName, tableName))
		r.con.Exec(fmt.Sprintf("ALTER TABLE %s RENAME TO %s", tableName+"_temp", simulationTableName))
	}
}
