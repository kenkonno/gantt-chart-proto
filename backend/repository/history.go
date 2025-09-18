package repository

import (
	"fmt"
	"strings"
	"time"

	"github.com/kenkonno/gantt-chart-proto/backend/models/db"
	"github.com/kenkonno/gantt-chart-proto/backend/repository/connection"
	"github.com/kenkonno/gantt-chart-proto/backend/repository/interfaces"
	"github.com/samber/lo"
	"gorm.io/gorm"
)

func NewHistoryRepository() interfaces.HistoryRepositoryIF {
	return &historyRepository{
		con: connection.GetCon(),
	}
}

type historyRepository struct {
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
	"units",
	"users",
	"facility_shared_links",
}

var columns = [][]string{
	{"id", "name", "color", "order", "created_at", "updated_at"},
	{"id", "name", "term_from", "term_to", "order", "status", "type", "shipment_due_date", "created_at", "updated_at"},
	{"id", "facility_id", "date", "type", "created_at", "updated_at"},
	{"id", "facility_id", "unit_id", "created_at", "updated_at"},
	{"id", "facility_id", "name", "date", "created_at", "updated_at"},
	{"id", "facility_id", "date", "description", "order", "created_at", "updated_at"},
	{"id", "facility_id", "unit_id", "process_id", "work_hour", "created_at", "updated_at"},
	{"id", "name", "order", "color", "created_at", "updated_at"},
	{"id", "gantt_group_id", "process_id", "department_id", "limit_date", "estimate", "number_of_worker", "days_after", "start_date", "end_date", "progress_percent", "memo", "order", "created_at", "updated_at"},
	{"id", "ticket_id", "user_id", "order", "created_at", "updated_at"},
	{"id", "name", "facility_id", "order", "created_at", "updated_at"},
	{"id", "department_id", "limit_of_operation", "last_name", "first_name", "password", "email", "role", "password_reset", "employment_start_date", "employment_end_date", "created_at", "updated_at"},
	{"id", "facility_id", "uuid", "created_at", "updated_at"},
}

func (r *historyRepository) CreateSnapshot(facilityId int32, name string) (int32, error) {
	tx := r.con.Begin()
	if tx.Error != nil {
		return 0, tx.Error
	}

	defer func() {
		if rec := recover(); rec != nil {
			tx.Rollback()
		}
	}()

	// 1. Create a new record in the facility_histories table.
	historyRecord := db.HistoryFacility{
		FacilityId: facilityId,
		Name:       name,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	if err := tx.Create(&historyRecord).Error; err != nil {
		tx.Rollback()
		return 0, err
	}
	historyId := *historyRecord.Id

	// 2. Copy data to history tables
	for i, tableName := range tableNames {
		historyTableName := "history_" + tableName
		cols := columns[i]

		quotedCols := lo.Map(cols, func(item string, index int) string {
			return `"` + item + `"`
		})

		insertCols := append([]string{`"history_id"`}, quotedCols...)

		q := fmt.Sprintf("INSERT INTO %s (%s) SELECT ?, %s FROM %s",
			historyTableName,
			strings.Join(insertCols, ", "),
			strings.Join(quotedCols, ", "),
			tableName)

		if err := tx.Exec(q, historyId).Error; err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	if err := tx.Commit().Error; err != nil {
		return 0, err
	}

	return historyId, nil
}

func (r *historyRepository) FindByFacilityId(facilityId int32) []db.HistoryFacility {
	var histories []db.HistoryFacility
	result := r.con.Where("facility_id = ?", facilityId).Order("created_at DESC").Find(&histories)
	if result.Error != nil {
		panic(result.Error)
	}
	return histories
}

func (r *historyRepository) UpdateName(id int32, name string) error {
	result := r.con.Model(&db.HistoryFacility{}).Where("id = ?", id).Update("name", name)
	return result.Error
}

func (r *historyRepository) Delete(id int32) error {
	tx := r.con.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	defer func() {
		if rec := recover(); rec != nil {
			tx.Rollback()
		}
	}()

	// Delete from all history tables
	for _, tableName := range tableNames {
		historyTableName := "history_" + tableName
		if err := tx.Exec(fmt.Sprintf("DELETE FROM %s WHERE history_id = ?", historyTableName), id).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	// Delete from the main history table
	if err := tx.Where("id = ?", id).Delete(&db.HistoryFacility{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
