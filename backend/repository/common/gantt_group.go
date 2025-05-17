package common

import (
	"fmt"
	"github.com/kenkonno/gantt-chart-proto/backend/models/db"
	"github.com/kenkonno/gantt-chart-proto/backend/repository/connection"
	"github.com/kenkonno/gantt-chart-proto/backend/repository/interfaces"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Auto generated start
func NewGanttGroupRepository() interfaces.GanttGroupRepositoryIF {
	return &ganttGroupRepository{connection.GetCon()}
}

type ganttGroupRepository struct {
	con *gorm.DB
}

func (r *ganttGroupRepository) FindAll() []db.GanttGroup {
	var ganttGroups []db.GanttGroup

	result := r.con.Order("id DESC").Find(&ganttGroups)
	if result.Error != nil {
		panic(result.Error)
	}
	return ganttGroups
}

func (r *ganttGroupRepository) Find(id int32) db.GanttGroup {
	var ganttGroup db.GanttGroup

	result := r.con.First(&ganttGroup, id)
	if result.Error != nil {
		panic(result.Error)
	}
	return ganttGroup
}

func (r *ganttGroupRepository) Upsert(m db.GanttGroup) db.GanttGroup {
	r.con.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		UpdateAll: true,
	}).Create(&m)
	return m
}

func (r *ganttGroupRepository) Delete(id int32) {
	r.con.Where("id = ?", id).Delete(db.GanttGroup{})
}

// Auto generated end
func (r *ganttGroupRepository) FindByFacilityId(facilityId []int32) []db.GanttGroup {
	var results []db.GanttGroup

	r.con.Raw(fmt.Sprintf(`
	SELECT
		gg.id
	,   gg.facility_id
	,   gg.unit_id
	FROM
		gantt_groups gg
	INNER JOIN
		units u 
	ON
		gg.unit_id = u.id
	WHERE
		gg.facility_id IN %s
	ORDER BY u.order
	`, connection.CreateInParamInt32(facilityId))).Scan(&results)

	return results
}

func (r *ganttGroupRepository) DeleteByUnitId(unitId int32) {
	r.con.Where("unit_id = ?", unitId).Delete(db.GanttGroup{})
}
