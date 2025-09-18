package history

import (
	"fmt"
	"github.com/kenkonno/gantt-chart-proto/backend/models/db"
	"github.com/kenkonno/gantt-chart-proto/backend/repository/connection"
	"github.com/kenkonno/gantt-chart-proto/backend/repository/interfaces"
	"gorm.io/gorm"
)

func NewGanttGroupRepository(historyId int32) interfaces.GanttGroupRepositoryIF {
	return &ganttGroupRepository{
		historyId: historyId,
		con:       connection.GetCon(),
	}
}

type ganttGroupRepository struct {
	historyId int32
	con       *gorm.DB
}

func (r *ganttGroupRepository) FindAll() []db.GanttGroup {
	var ganttGroups []db.GanttGroup
	result := r.con.Table("history_gantt_groups").Where("history_id = ?", r.historyId).Order("id DESC").Find(&ganttGroups)
	if result.Error != nil {
		panic(result.Error)
	}
	return ganttGroups
}

func (r *ganttGroupRepository) Find(id int32) db.GanttGroup {
	var ganttGroup db.GanttGroup
	result := r.con.Table("history_gantt_groups").Where("history_id = ? AND id = ?", r.historyId, id).First(&ganttGroup)
	if result.Error != nil {
		panic(result.Error)
	}
	return ganttGroup
}

func (r *ganttGroupRepository) Upsert(m db.GanttGroup) db.GanttGroup {
	// History is read-only
	return m
}

func (r *ganttGroupRepository) Delete(id int32) {
	// History is read-only
}

func (r *ganttGroupRepository) FindByFacilityId(facilityId []int32) []db.GanttGroup {
	var results []db.GanttGroup

	r.con.Raw(fmt.Sprintf(`
	SELECT
		gg.id
	,   gg.facility_id
	,   gg.unit_id
	FROM
		history_gantt_groups gg
	INNER JOIN
		history_units u
	ON
		gg.unit_id = u.id
	WHERE
		gg.history_id = ? AND u.history_id = ? AND gg.facility_id IN ?
	ORDER BY u.order
	`, connection.CreateInParamInt32(facilityId)), r.historyId, r.historyId, facilityId).Scan(&results)

	return results
}

func (r *ganttGroupRepository) DeleteByUnitId(unitId int32) {
	// History is read-only
}
