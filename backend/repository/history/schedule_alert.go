package history

import (
	"github.com/kenkonno/gantt-chart-proto/backend/models/db"
	"github.com/kenkonno/gantt-chart-proto/backend/repository/connection"
	"github.com/kenkonno/gantt-chart-proto/backend/repository/interfaces"
	"gorm.io/gorm"
)

func NewScheduleAlertRepository(historyId int32) interfaces.ScheduleAlertRepositoryIF {
	return &scheduleAlertRepository{
		historyId: historyId,
		con:       connection.GetCon(),
	}
}

type scheduleAlertRepository struct {
	historyId int32
	con       *gorm.DB
}

func (r *scheduleAlertRepository) FindAll() []db.ScheduleAlert {
	var results []db.ScheduleAlert

	// Note: This is a complex raw query adapted from the simulation repository.
	// It has been modified to use history tables and filter by history_id.
	r.con.Raw(`
	-- 稼働設定は初期値がない場合があるのでクエリーで生成する
	WITH all_operation_settings AS (
		SELECT
			os.id
			 ,   u.id as unit_id
			 ,   p.id as process_id
			 ,   COALESCE(os.work_hour, 8) as work_hour
			 ,   COALESCE(os.created_at, now()) as created_at
			 ,   os.updated_at
			 , os.facility_id
		FROM
			history_processes p
				CROSS JOIN
			history_units u
				LEFT JOIN
			history_operation_settings os
			ON
						os.unit_id = u.id
					AND os.process_id = p.id
					AND os.history_id = ?
		WHERE u.history_id = ? AND p.history_id = ?
		ORDER BY u.id, p.id
	), target_tickets AS (
		SELECT
			f.id as facility_id
			 ,   f.name as facility_name
			 ,   u.id as unit_id
			 ,   u.name as unit_name
			 ,   p.id as process_id
			 ,   p.name as process_name
			 ,   COALESCE(os.work_hour, 8) as work_hour
			 ,   COALESCE(t.estimate, 0) as estimate
			 ,   t.number_of_worker
			 ,   COALESCE(t.progress_percent, 0) as progress_percent
			 ,   t.start_date
			 ,   t.end_date
			 ,   t.id as ticket_id
			 ,   (SELECT ARRAY_AGG(date ORDER BY date)
				  FROM generate_series(t.start_date, f.term_to, interval '1days') as date
				  WHERE NOT EXISTS (SELECT * FROM (SELECT date FROM history_holidays
					  UNION
					  SELECT date FROM history_facility_work_schedules fws WHERE type = 'Holiday' AND fws.facility_id = f.id
					  EXCEPT
					  SELECT date FROM history_facility_work_schedules fws WHERE type = 'WorkingDay' AND fws.facility_id = f.id) as h WHERE date.date = h.date)
				  ) as term_business_dates
			 ,   (SELECT ARRAY_AGG(date ORDER BY date)
				  FROM generate_series(t.start_date, t.end_date, interval '1days') as date
				  WHERE NOT EXISTS (SELECT * FROM (SELECT date FROM history_holidays
					  UNION
					  SELECT date FROM history_facility_work_schedules fws WHERE type = 'Holiday' AND fws.facility_id = f.id
					  EXCEPT
					  SELECT date FROM history_facility_work_schedules fws WHERE type = 'WorkingDay' AND fws.facility_id = f.id) as h WHERE date.date = h.date)
				) as ticket_business_dates
			 ,   f."order" as facility_order
			 ,   p."order" as process_order
			 ,   t."order" as ticket_order
		FROM history_facilities f
				 INNER JOIN history_gantt_groups gg on f.id = gg.facility_id AND gg.history_id = f.history_id
				 INNER JOIN history_units u on u.id = gg.unit_id AND u.history_id = f.history_id
				 INNER JOIN history_tickets t on gg.id = t.gantt_group_id AND t.history_id = f.history_id
				 INNER JOIN history_processes p on t.process_id = p.id AND p.history_id = f.history_id
				 LEFT JOIN all_operation_settings os ON os.facility_id = f.id AND os.unit_id = gg.unit_id AND os.process_id = p.id
		WHERE (t.progress_percent <> 100 OR t.progress_percent IS NULL)
		  AND t.start_date IS NOT NULL
		  AND t.end_date IS NOT NULL
		  AND t.process_id IS NOT NULL
		  AND t.start_date < NOW() -- まだ開始していないもの除外
		  AND f.status = 'Enabled'
          AND f.type = 'Ordered'
		  AND f.history_id = ?
	)
	-- 対象のチケットから計算を行う
	-- 開始日から終了日の内休日を取り除く。
	-- 稼働予定日数（営業日ベース）を進捗で割った箇所が現在の日付となる。
	SELECT
		facility_id
		 ,  facility_name
		 ,  unit_id
		 ,  unit_name
		 ,  process_id
		 ,  process_name
		 ,  start_date
		 ,  end_date
		 ,  progress_percent
		 , ticket_id
		 , work_hour
		 , number_of_worker
		 , ticket_business_dates[array_length(ticket_business_dates, 1) * progress_percent / 100] as actual_progress_date
		 , (SELECT COUNT(dates)
			FROM UNNEST(term_business_dates[array_length(ticket_business_dates, 1) * progress_percent / 100:]) as dates
			WHERE dates < NOW()
			) - 2 + CASE WHEN array_length(ticket_business_dates, 1) * progress_percent / 100 = 0 THEN 1 ELSE 0 END as delay_days -- 当日と進捗による完了日を除く,1日も進捗がない場合は遅延を+1する
	FROM
		target_tickets
	WHERE
		(SELECT COUNT(dates)
		 FROM UNNEST(term_business_dates[array_length(ticket_business_dates, 1) * progress_percent / 100:]) as dates
		 WHERE dates < NOW()
		) - 2 + CASE WHEN array_length(ticket_business_dates, 1) * progress_percent / 100 = 0 THEN 1 ELSE 0 END > 0
	ORDER BY
		facility_order
		,   ticket_order
		`, r.historyId, r.historyId, r.historyId, r.historyId).Scan(&results)
	return results
}

func (r *scheduleAlertRepository) Find(id int32) db.ScheduleAlert {
	// This model does not have a primary key in the database, so finding by ID is not applicable.
	return db.ScheduleAlert{}
}
