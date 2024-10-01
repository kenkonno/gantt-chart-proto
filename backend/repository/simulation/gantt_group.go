package simulation

import (
	"fmt"
	"github.com/kenkonno/gantt-chart-proto/backend/models/db"
	"github.com/kenkonno/gantt-chart-proto/backend/repository/connection"
	"github.com/kenkonno/gantt-chart-proto/backend/repository/interfaces"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Auto generated start
func NewSimulationGanttGroupRepository() interfaces.GanttGroupRepositoryIF {
	return &ganttGroupRepository{
		con:   connection.GetCon(),
		table: "simulation_gantt_groups",
	}
}

type ganttGroupRepository struct {
	con *gorm.DB
	table string
}

func (r *ganttGroupRepository) FindAll() []db.GanttGroup {
	var ganttGroups []db.GanttGroup

	result := r.con.Table(r.table).Order("id DESC").Find(&ganttGroups)
	if result.Error != nil {
		panic(result.Error)
	}
	return ganttGroups
}

func (r *ganttGroupRepository) Find(id int32) db.GanttGroup {
	var ganttGroup db.GanttGroup

	result := r.con.Table(r.table).First(&ganttGroup, id)
	if result.Error != nil {
		panic(result.Error)
	}
	return ganttGroup
}

func (r *ganttGroupRepository) Upsert(m db.GanttGroup) db.GanttGroup {
	r.con.Table(r.table).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		UpdateAll: true,
	}).Create(&m)
	return m
}

func (r *ganttGroupRepository) Delete(id int32) {
	r.con.Table(r.table).Where("id = ?", id).Delete(db.GanttGroup{})
}

// Auto generated end
func (r *ganttGroupRepository) FindByFacilityId(facilityId int32) []db.GanttGroup {
	var results []db.GanttGroup

	r.con.Table(r.table).Raw(fmt.Sprintf(`
	SELECT
		gg.id
	,   %d as facility_id
	,   gg.unit_id
	FROM
		simulation_gantt_groups gg
	INNER JOIN
		simulation_units u 
	ON
		gg.unit_id = u.id
	WHERE
		gg.facility_id = %d
	ORDER BY u.order
	`, facilityId, facilityId)).Scan(&results)

	return results
}

func (r *ganttGroupRepository) DeleteByUnitId(unitId int32) {
	r.con.Table(r.table).Where("unit_id = ?", unitId).Delete(db.GanttGroup{})
}
