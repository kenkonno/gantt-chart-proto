package simulation

import (
	"fmt"
	"github.com/kenkonno/gantt-chart-proto/backend/models/db"
	"github.com/kenkonno/gantt-chart-proto/backend/repository/connection"
	"github.com/kenkonno/gantt-chart-proto/backend/repository/interfaces"
	"gorm.io/gorm"
)

// Auto generated start
func NewSimulationPileUpsRepository() interfaces.PileUpsRepositoryIF {
	return &pileUpsRepository{
		con:   connection.GetCon(),
		table: "simulation_pileUps",
	}
}

type pileUpsRepository struct {
	con   *gorm.DB
	table string
}

// GetDefaultPileUps 全ての設備の積み上げを返却する
func (r *pileUpsRepository) GetDefaultPileUps(excludeFacilityId int32, facilityTypes []string) []db.DefaultPileUp {
	var results []db.DefaultPileUp

	r.con.Table(r.table).Raw(fmt.Sprintf(`
	WITH w_facilities AS (
		SELECT MIN(term_from) as min_date, MAX(term_to) as max_date FROM simulation_facilities WHERE type IN %s AND status IN ('Enabled')
	), date_master AS (
		SELECT
			date.date
			 , ARRAY_REMOVE(ARRAY_AGG(DISTINCT th.facility_id), NULL) as facility_ids -- 祝日の設備一覧
			 , ARRAY_REMOVE(ARRAY_AGG(DISTINCT u.id), NULL) as user_ids -- 稼動可能なユーザー一覧
			 , ROW_NUMBER() OVER(ORDER BY date.date) - 1 as index
		FROM generate_series((SELECT min_date FROM w_facilities), (SELECT max_date FROM w_facilities), interval '1days') as date
				 LEFT JOIN simulation_holidays th ON date.date = th.date
				 LEFT JOIN simulation_users u ON (u.employment_start_date <= date.date AND (date.date <= u.employment_end_date OR u.employment_end_date IS NULL) )
		GROUP BY
			date.date
	), target_tickets_by_user AS (
		SELECT
			t.id
			 ,   f.id as facility_id
			 ,   t.department_id
			 ,   t.estimate
			 ,   t.number_of_worker
			 ,   t.start_date
			 ,   t.end_date
			 ,   tu.user_id
			 ,   (SELECT ARRAY_AGG(index) FROM date_master dm WHERE dm.date BETWEEN t.start_date AND t.end_date AND NOT(f.id = ANY (dm.facility_ids))) as valid_indexes
		FROM simulation_tickets t
				 INNER JOIN simulation_gantt_groups gg ON gg.id = t.gantt_group_id
				 INNER JOIN simulation_facilities f ON f.id = gg.facility_id AND f.id != %d AND type IN %s AND status IN ('Enabled')
				 LEFT JOIN simulation_ticket_users tu ON t.id = tu.ticket_id
		WHERE
			gantt_group_id IN (SELECT id FROM simulation_gantt_groups WHERE facility_id = f.id)
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
			   ,   tu.user_id
	)
	SELECT
        ttbu.id
         ,   ttbu.facility_id
         ,   ttbu.department_id
         ,   ttbu.estimate
         ,   ttbu.number_of_worker
         ,   ttbu.start_date
         ,   ttbu.end_date
		 ,   ARRAY_REMOVE(ARRAY_AGG(DISTINCT ttbu.user_id), null) as user_ids
         ,   ttbu.valid_indexes
         ,   COUNT(*) as number_of_worker_by_day
	     ,   (SELECT COUNT(*) FROM date_master dm WHERE dm.date BETWEEN ttbu.start_date AND ttbu.end_date AND NOT( ttbu.facility_id = ANY (dm.facility_ids))) as number_of_work_day
    FROM
        target_tickets_by_user ttbu
    LEFT JOIN
        date_master dm ON dm.index = ANY(ttbu.valid_indexes)AND ttbu.user_id = ANY(dm.user_ids)
    GROUP BY
        ttbu.id
         ,   ttbu.facility_id
         ,   ttbu.department_id
         ,   ttbu.estimate
         ,   ttbu.number_of_worker
         ,   ttbu.start_date
         ,   ttbu.end_date
         ,   ttbu.valid_indexes
	`, connection.CreateInParam(facilityTypes), excludeFacilityId, connection.CreateInParam(facilityTypes))).Scan(&results)

	return results
}

// GetUserInfos validIndex毎にどのユーザーが稼動可能なのかを格納した情報
func (r *pileUpsRepository) GetValidIndexUsers() []db.ValidIndexUser {
	var results []db.ValidIndexUser

	r.con.Raw(fmt.Sprintf(`
	WITH w_facilities AS (
		SELECT MIN(term_from) as min_date, MAX(term_to) as max_date FROM simulation_facilities WHERE status IN ('Enabled')
	)
	SELECT
		date.date
		 , ARRAY_REMOVE(ARRAY_AGG(DISTINCT th.facility_id), NULL) as facility_ids -- 祝日の設備一覧
		 , ARRAY_REMOVE(ARRAY_AGG(DISTINCT u.id), NULL) as user_ids -- 稼動可能なユーザー一覧
		 , ROW_NUMBER() OVER(ORDER BY date.date) - 1 as index
	FROM generate_series((SELECT min_date FROM w_facilities), (SELECT max_date FROM w_facilities), interval '1days') as date
			 LEFT JOIN simulation_holidays th ON date.date = th.date
			 LEFT JOIN simulation_users u ON (u.employment_start_date <= date.date AND (date.date <= u.employment_end_date OR u.employment_end_date IS NULL) )
	GROUP BY
		date.date
	`)).Scan(&results)

	return results
}
