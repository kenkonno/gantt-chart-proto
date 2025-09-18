package simulation

import (
	"fmt"
	"strings"

	"github.com/kenkonno/gantt-chart-proto/backend/repository/connection"
	"github.com/samber/lo"
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

var tableNames = []string{
	"departments",
	"facilities",
	"facility_work_schedules",
	"gantt_groups",
	"holidays",
	"milestones",
	"operation_settings",
	"processes",
	"tickets",
	"ticket_users",
	"ticket_daily_weights",
	"units",
	"users",
}

var columns = [][]string{
	{"id", "name", "order", "created_at", "updated_at"},
	{"id", "name", "term_from", "term_to", "order", "created_at", "updated_at", "status", "type", "shipment_due_date", "free_text"},
	{"id", "facility_id", "date", "type", "created_at", "updated_at"},
	{"id", "facility_id", "unit_id", "created_at", "updated_at"},
	{"id", "facility_id", "name", "date", "created_at", "updated_at"},
	{"id", "facility_id", "date", "description", "order", "created_at", "updated_at"},
	{"id", "facility_id", "unit_id", "process_id", "work_hour", "created_at", "updated_at"},
	{"id", "name", "order", "created_at", "updated_at", "color"},
	{"id", "gantt_group_id", "process_id", "department_id", "limit_date", "estimate", "number_of_worker", "days_after", "start_date", "end_date", "progress_percent", "order", "created_at", "updated_at", "memo"},
	{"id", "ticket_id", "user_id", "order", "created_at", "updated_at"},
	{"ticket_id", "work_hour", "date", "created_at", "updated_at"},
	{"id", "name", "facility_id", "order", "created_at", "updated_at"},
	{"id", "department_id", "limit_of_operation", "password", "email", "created_at", "updated_at", "role", "last_name", "first_name", "password_reset", "employment_start_date", "employment_end_date"},
}

func (r *simulationRepository) InitAllData() {

	// Truncate, Copy
	for index, tableName := range tableNames {
		column := strings.Join(lo.Map(columns[index], func(item string, index int) string {
			return `"` + item + `"`
		}), ",")
		r.con.Exec(fmt.Sprintf("TRUNCATE TABLE %s", "simulation_"+tableName))
		r.con.Exec(fmt.Sprintf("INSERT INTO %s (%s) SELECT %s FROM %s", "simulation_"+tableName, column, column, tableName))
	}
}

func (r *simulationRepository) ResetSequence(con *gorm.DB) {

	var db = r.con
	if con != nil {
		db = con
	}

	// Truncate, Copy
	for _, tableName := range tableNames {
		simulationTableName := "simulation_" + tableName
		db.Exec(fmt.Sprintf("SELECT setval('%s', (SELECT MAX(id) FROM %s) + 1)", tableName+"_id_seq", tableName))
		db.Exec(fmt.Sprintf("SELECT setval('%s', (SELECT MAX(id) FROM %s) + 1)", simulationTableName+"_id_seq", simulationTableName))
	}
}

func (r *simulationRepository) SwitchTable() {

	// トランザクション開始
	tx := r.con.Begin()
	if tx.Error != nil {
		panic(tx.Error)
	}

	// エラー時のロールバック処理
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	// Truncate, Copy
	// INSERTでシミュレーションテーブルから本番テーブルにデータをクリーンコピー
	for index, tableName := range tableNames {
		column := strings.Join(lo.Map(columns[index], func(item string, index int) string {
			return `"` + item + `"`
		}), ",")
		simulationTableName := "simulation_" + tableName

		// 本番テーブルをクリアしてシミュレーションテーブルのデータを挿入
		result := tx.Exec(fmt.Sprintf("TRUNCATE TABLE %s", tableName))
		if result.Error != nil {
			tx.Rollback()
			panic(result.Error)
		}

		result = tx.Exec(fmt.Sprintf("INSERT INTO %s (%s) SELECT %s FROM %s", tableName, column, column, simulationTableName))
		if result.Error != nil {
			tx.Rollback()
			panic(result.Error)
		}
	}
	// トランザクションコミット
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		panic(err)
	}
}
