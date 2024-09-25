package common

import (
	"fmt"
	"github.com/kenkonno/gantt-chart-proto/backend/models/db"
	"github.com/kenkonno/gantt-chart-proto/backend/repository/connection"
	"github.com/kenkonno/gantt-chart-proto/backend/repository/interfaces"
	"gorm.io/gorm"
)

// Auto generated start
func NewPileUpsRepository() interfaces.PileUpsRepositoryIF {
	return &pileUpsRepository{connection.GetCon()}
}

type pileUpsRepository struct {
	con *gorm.DB
}

// GetDefaultPileUps 全ての設備の積み上げを返却する
func (r *pileUpsRepository) GetDefaultPileUps(excludeFacilityId int32, facilityTypes []string) []db.DefaultPileUp {
	var results []db.DefaultPileUp

	r.con.Raw(fmt.Sprintf(`
	WITH w_facilities AS (
		SELECT MIN(term_from) as min_date, MAX(term_to) as max_date FROM facilities WHERE id != %d AND type IN %s AND status IN ('Enabled')
	), date_master AS (
		SELECT
			date.date
			 , ARRAY_REMOVE(ARRAY_AGG(th.facility_id), NULL) as facility_ids
			 , ROW_NUMBER() OVER(ORDER BY date.date) - 1 as index
		FROM generate_series((SELECT min_date FROM w_facilities), (SELECT max_date FROM w_facilities), interval '1days') as date
				 LEFT JOIN holidays th ON date.date = th.date
		GROUP BY
			date.date
	), target_tickets AS (
		SELECT
			t.id
			 ,   f.id as facility_id
			 ,   t.department_id
			 ,   t.estimate
			 ,   t.number_of_worker
			 ,   t.start_date
			 ,   t.end_date
			 ,   ARRAY_REMOVE(ARRAY_AGG(tu.user_id), null) as user_ids
			 ,   t.estimate::numeric / (SELECT COUNT(*) FROM date_master dm WHERE dm.date BETWEEN t.start_date AND t.end_date AND NOT(f.id = ANY (dm.facility_ids))) as work_per_day
			 ,   (SELECT ARRAY_AGG(index) FROM date_master dm WHERE dm.date BETWEEN t.start_date AND t.end_date AND NOT(f.id = ANY (dm.facility_ids))) as valid_indexes
		FROM tickets t
				 INNER JOIN gantt_groups gg ON gg.id = t.gantt_group_id
				 INNER JOIN facilities f ON f.id = gg.facility_id AND f.id != %d AND type IN %s AND status IN ('Enabled')
				 LEFT JOIN ticket_users tu ON t.id = tu.ticket_id
		WHERE
			gantt_group_id IN (SELECT id FROM gantt_groups WHERE facility_id = f.id)
		  AND t.number_of_worker > 0
		  AND (SELECT COUNT(*) FROM date_master dm WHERE dm.date BETWEEN t.start_date AND t.end_date AND NOT(f.id = ANY (dm.facility_ids))) > 0
		  AND estimate > 0
		GROUP BY
			t.id
			   ,   f.id
			   ,   t.department_id
			   ,   t.estimate
			   ,   t.number_of_worker
			   ,   t.start_date
			   ,   t.end_date
	)
	SELECT * FROM target_tickets tt
	`, excludeFacilityId, connection.CreateInParam(facilityTypes), excludeFacilityId, connection.CreateInParam(facilityTypes))).Scan(&results)

	return results
}
