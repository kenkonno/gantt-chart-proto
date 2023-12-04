package repository

import (
	"github.com/kenkonno/gantt-chart-proto/backend/models/db"
	"gorm.io/gorm"
)

// Auto generated start
func NewScheduleAlertRepository() scheduleAlertRepository {
	return scheduleAlertRepository{con}
}

type scheduleAlertRepository struct {
	con *gorm.DB
}

func (r *scheduleAlertRepository) FindAll() []db.ScheduleAlert {
	var results []db.ScheduleAlert

	r.con.Raw(`
	-- 稼働設定は初期値がない場合があるのでクエリーで生成する
	WITH all_operation_settings AS (
		SELECT
			operation_settings.id
			 ,   units.id as unit_id
			 ,   processes.id as process_id
			 ,   COALESCE(operation_settings.work_hour, 8) as work_hour
			 ,   COALESCE(operation_settings.created_at, now()) as created_at
			 ,   operation_settings.updated_at
			 , operation_settings.facility_id
		FROM
			processes
				CROSS JOIN
			units
				LEFT JOIN
			operation_settings
			ON
						operation_settings.unit_id = units.id
					AND operation_settings.process_id = processes.id
		ORDER BY units.id, processes.id
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
				  WHERE date NOT IN (SELECT h.date FROM holidays h WHERE f.id = h.facility_id)
				  ) as term_business_dates
			 ,   (SELECT ARRAY_AGG(date ORDER BY date)
				  FROM generate_series(t.start_date, t.end_date, interval '1days') as date
				  WHERE date NOT IN (SELECT h.date FROM holidays h WHERE f.id = h.facility_id)
				) as ticket_business_dates
			 ,   f."order" as facility_order
			 ,   p."order" as process_order
			 ,   t."order" as ticket_order
		FROM facilities f
				 INNER JOIN gantt_groups gg on f.id = gg.facility_id
				 INNER JOIN units u on u.id = gg.unit_id
				 INNER JOIN tickets t on gg.id = t.gantt_group_id
				 INNER JOIN processes p on t.process_id = p.id
				 LEFT JOIN all_operation_settings os ON os.facility_id = f.id AND os.unit_id = gg.unit_id AND os.process_id = p.id
		WHERE (t.progress_percent <> 100 OR t.progress_percent IS NULL)
		  AND t.start_date IS NOT NULL
		  AND t.end_date IS NOT NULL
		  AND t.process_id IS NOT NULL
		  AND t.start_date < NOW() -- まだ開始していないもの除外
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
		`).Scan(&results)
	return results

}

func (r *scheduleAlertRepository) Find(id int32) db.ScheduleAlert {
	var scheduleAlert db.ScheduleAlert

	result := r.con.First(&scheduleAlert, id)
	if result.Error != nil {
		panic(result.Error)
	}
	return scheduleAlert
}
