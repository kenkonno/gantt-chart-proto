package common

import (
	"fmt"
	"github.com/kenkonno/gantt-chart-proto/backend/models/db"
	"github.com/kenkonno/gantt-chart-proto/backend/repository/connection"
	"github.com/kenkonno/gantt-chart-proto/backend/repository/interfaces"
	"gorm.io/gorm"
	"time"
)

// Auto generated start
func NewPileUpsRepository() interfaces.PileUpsRepositoryIF {
	return &pileUpsRepository{connection.GetCon()}
}

type pileUpsRepository struct {
	con *gorm.DB
}

// GetDefaultPileUps 全ての設備の積み上げを返却する。クエリーが非常に複雑なのでメモを残す。
// w_facilitiesは自分以外の設備を除外する
// date_masterは全設備の開始日からindexを採番している。
// 日付毎にその設備が祝日なのか、その日に稼動できるユーザーの一覧を持っている。
// target_tickets_by_userではチケット×紐づくユーザー分の行数となる。在籍期間対応でそうなった。どうしても日毎に稼動可能なユーザー数を計算することができなかったため。
// 最後のSELECTはチケット単位に日毎に稼動可能なユーザーを集計している。
func (r *pileUpsRepository) GetDefaultPileUps(excludeFacilityId int32, facilityTypes []string) []db.DefaultPileUp {
	var results []db.DefaultPileUp

	r.con.Debug().Raw(fmt.Sprintf(`
	WITH w_facilities AS (
		SELECT MIN(term_from) as min_date, MAX(term_to) as max_date FROM facilities WHERE type IN %s AND status IN ('Enabled')
	), date_master AS (
		SELECT
			date.date
			 , ARRAY_REMOVE(ARRAY_AGG(DISTINCT CASE WHEN tws.type = 'Holiday' THEN tws.facility_id END ), NULL) as holiday_facility_ids -- 祝日の設備一覧
			 , ARRAY_REMOVE(ARRAY_AGG(DISTINCT CASE WHEN tws.type = 'WorkingDay' THEN tws.facility_id END ), NULL) as working_day_facility_ids -- 稼働日の設備一覧
	         , th.date IS NOT NULL as is_holiday_by_master
			 , ARRAY_REMOVE(ARRAY_AGG(DISTINCT u.id), NULL) as user_ids -- 稼動可能なユーザー一覧
			 , ROW_NUMBER() OVER(ORDER BY date.date) - 1 as index
		FROM generate_series((SELECT min_date FROM w_facilities), (SELECT max_date FROM w_facilities), interval '1days') as date
				 LEFT JOIN holidays th ON date.date = th.date
	             LEFT JOIN facility_work_schedules tws ON date.date = tws.date
				 LEFT JOIN users u ON (u.employment_start_date <= date.date AND (date.date <= u.employment_end_date OR u.employment_end_date IS NULL) AND u.role IN ('worker','manager','viewer'))
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
                -- 稼働日として明示されている場合
                f.id = ANY(dm.working_day_facility_ids)
                OR
                -- または祝日でない場合
                NOT(
                    f.id = ANY(dm.holiday_facility_ids) -- 設備ごとの祝日
                        OR dm.is_holiday_by_master      -- マスター上の祝日
                    )
                )) as valid_indexes
		FROM tickets t
				 INNER JOIN gantt_groups gg ON gg.id = t.gantt_group_id
				 INNER JOIN facilities f ON f.id = gg.facility_id AND f.id != %d AND type IN %s AND status IN ('Enabled')
				 LEFT JOIN ticket_users tu ON t.id = tu.ticket_id
		WHERE
			gantt_group_id IN (SELECT id FROM gantt_groups WHERE facility_id = f.id)
		  AND t.number_of_worker > 0
		  AND (SELECT COUNT(*) FROM date_master dm WHERE dm.date BETWEEN t.start_date AND t.end_date AND (
                -- 稼働日として明示されている場合
                f.id = ANY(dm.working_day_facility_ids)
                OR
                -- または祝日でない場合
                NOT(
                    f.id = ANY(dm.holiday_facility_ids) -- 設備ごとの祝日
                        OR dm.is_holiday_by_master      -- マスター上の祝日
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
        -- 稼働日として明示されている場合
            ttbu.facility_id = ANY(dm.working_day_facility_ids)
        OR
        -- または祝日でない場合
        NOT(
            ttbu.facility_id = ANY(dm.holiday_facility_ids) -- 設備ごとの祝日
                OR dm.is_holiday_by_master      -- マスター上の祝日
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
	`, connection.CreateInParam(facilityTypes), excludeFacilityId, connection.CreateInParam(facilityTypes))).Scan(&results)

	return results
}

// GetUserInfos validIndex毎にどのユーザーが稼動可能なのかを格納した情報 祝日に関する処理はないので祝日の変更で修正なし。
func (r *pileUpsRepository) GetValidIndexUsers(globalStartDate time.Time) []db.ValidIndexUser {
	var results []db.ValidIndexUser

	r.con.Raw(fmt.Sprintf(`
	WITH w_facilities AS (
		SELECT MIN(term_from) as min_date, MAX(term_to) as max_date FROM facilities WHERE status IN ('Enabled')
	)
	SELECT
		date.date
		 , ARRAY_REMOVE(ARRAY_AGG(DISTINCT th.facility_id), NULL) as facility_ids -- 祝日の設備一覧
		 , ARRAY_REMOVE(ARRAY_AGG(DISTINCT u.id), NULL) as user_ids -- 稼動可能なユーザー一覧
		 , ROW_NUMBER() OVER(ORDER BY date.date) - 1 as valid_index
	FROM generate_series((SELECT min_date FROM w_facilities), (SELECT max_date FROM w_facilities), interval '1days') as date
			 LEFT JOIN holidays th ON date.date = th.date
			 LEFT JOIN users u ON (u.employment_start_date <= date.date AND (date.date <= u.employment_end_date OR u.employment_end_date IS NULL) AND u.role IN ('worker','manager','viewer'))
	WHERE date.date >= $1
	GROUP BY
		date.date
	`), globalStartDate).Scan(&results)

	return results
}
