package history

import (
	"fmt"
	"github.com/kenkonno/gantt-chart-proto/backend/models/db"
	"github.com/kenkonno/gantt-chart-proto/backend/repository/connection"
	"github.com/kenkonno/gantt-chart-proto/backend/repository/interfaces"
	"gorm.io/gorm"
	"time"
)

func NewPileUpsRepository(historyId int32) interfaces.PileUpsRepositoryIF {
	return &pileUpsRepository{
		historyId: historyId,
		con:       connection.GetCon(),
	}
}

type pileUpsRepository struct {
	historyId int32
	con       *gorm.DB
}

func (r *pileUpsRepository) GetDefaultPileUps(excludeFacilityId int32, facilityTypes []string) []db.DefaultPileUp {
	var results []db.DefaultPileUp

	r.con.Raw(fmt.Sprintf(`
	WITH w_facilities AS (
		SELECT MIN(term_from) as min_date, MAX(term_to) as max_date FROM history_facilities WHERE history_id = ? AND type IN %s AND status IN ('Enabled')
	), date_master AS (
		SELECT
			date.date
			 , ARRAY_REMOVE(ARRAY_AGG(DISTINCT CASE WHEN tws.type = 'Holiday' THEN tws.facility_id END ), NULL) as holiday_facility_ids
			 , ARRAY_REMOVE(ARRAY_AGG(DISTINCT CASE WHEN tws.type = 'WorkingDay' THEN tws.facility_id END ), NULL) as working_day_facility_ids
	         , th.date IS NOT NULL as is_holiday_by_master
			 , ARRAY_REMOVE(ARRAY_AGG(DISTINCT u.id), NULL) as user_ids
			 , ROW_NUMBER() OVER(ORDER BY date.date) - 1 as index
		FROM generate_series((SELECT min_date FROM w_facilities), (SELECT max_date FROM w_facilities), interval '1days') as date
				 LEFT JOIN history_holidays th ON date.date = th.date AND th.history_id = ?
	             LEFT JOIN history_facility_work_schedules tws ON date.date = tws.date AND tws.history_id = ?
				 LEFT JOIN history_users u ON (u.employment_start_date <= date.date AND (date.date <= u.employment_end_date OR u.employment_end_date IS NULL) AND u.role IN ('worker','manager','viewer')) AND u.history_id = ?
		GROUP BY
	        date.date, th.date
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
			 ,   (SELECT ARRAY_AGG(index) FROM date_master dm WHERE dm.date BETWEEN t.start_date AND t.end_date AND (
                f.id = ANY(dm.working_day_facility_ids)
                OR
                NOT(
                    f.id = ANY(dm.holiday_facility_ids)
                        OR dm.is_holiday_by_master
                    )
                )) as valid_indexes
		FROM history_tickets t
				 INNER JOIN history_gantt_groups gg ON gg.id = t.gantt_group_id AND gg.history_id = t.history_id
				 INNER JOIN history_facilities f ON f.id = gg.facility_id AND f.history_id = t.history_id AND f.id != %d AND type IN %s AND status IN ('Enabled')
				 LEFT JOIN history_ticket_users tu ON t.id = tu.ticket_id AND tu.history_id = t.history_id
		WHERE
			t.history_id = ?
		  AND t.number_of_worker > 0
		  AND (SELECT COUNT(*) FROM date_master dm WHERE dm.date BETWEEN t.start_date AND t.end_date AND (
                f.id = ANY(dm.working_day_facility_ids)
                OR
                NOT(
                    f.id = ANY(dm.holiday_facility_ids)
                        OR dm.is_holiday_by_master
                    )
                )) > 0
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
	     ,   (SELECT COUNT(*) FROM date_master dm WHERE dm.date BETWEEN ttbu.start_date AND ttbu.end_date AND (
            ttbu.facility_id = ANY(dm.working_day_facility_ids)
        OR
        NOT(
            ttbu.facility_id = ANY(dm.holiday_facility_ids)
                OR dm.is_holiday_by_master
            )
    )) as number_of_work_day
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
	`, connection.CreateInParam(facilityTypes), excludeFacilityId, connection.CreateInParam(facilityTypes)), r.historyId, r.historyId, r.historyId, r.historyId, r.historyId).Scan(&results)

	return results
}

func (r *pileUpsRepository) GetValidIndexUsers(globalStartDate time.Time) []db.ValidIndexUser {
	var results []db.ValidIndexUser

	r.con.Raw(fmt.Sprintf(`
	WITH w_facilities AS (
		SELECT MIN(term_from) as min_date, MAX(term_to) as max_date FROM history_facilities WHERE history_id = ? AND status IN ('Enabled')
	)
	SELECT
		date.date
		 , MAX(th.id) IS NOT NULL as is_holiday
		 , ARRAY_REMOVE(ARRAY_AGG(DISTINCT u.id), NULL) as user_ids
		 , ROW_NUMBER() OVER(ORDER BY date.date) - 1 as valid_index
	FROM generate_series((SELECT min_date FROM w_facilities), (SELECT max_date FROM w_facilities), interval '1days') as date
			 LEFT JOIN history_holidays th ON date.date = th.date AND th.history_id = ?
			 LEFT JOIN history_users u ON (u.employment_start_date <= date.date AND (date.date <= u.employment_end_date OR u.employment_end_date IS NULL) AND u.role IN ('worker','manager','viewer')) AND u.history_id = ?
	WHERE date.date >= $1
	GROUP BY
		date.date
	`), r.historyId, r.historyId, r.historyId, globalStartDate).Scan(&results)

	return results
}
