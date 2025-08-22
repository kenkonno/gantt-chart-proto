package default_pile_ups

import (
	"github.com/kenkonno/gantt-chart-proto/backend/api/constants"
	"github.com/kenkonno/gantt-chart-proto/backend/api/utils"
	"github.com/kenkonno/gantt-chart-proto/backend/models/db"
	"github.com/kenkonno/gantt-chart-proto/backend/repository/interfaces"
	"github.com/samber/lo"
	"time"
)

type WorkPerDayCalculator interface {
	CalculateWorkPerDay(ticket *db.DefaultPileUp, validUserMap map[int32][]int32) float32
	ReCalculateWorkPerDayByValidIndex(workHour float32, validIndex int32, ticket *db.DefaultPileUp) float32
}

func NewWorkPerDayCalculatorFactory(
	ticketDailyWeightRep interfaces.TicketDailyWeightRepositoryIF,
	globalStartDate time.Time,
	ticket db.DefaultPileUp,
	validUserMap map[int32][]int32,
) WorkPerDayCalculator {
	if utils.HasOption(constants.WorkloadWeighting) {
		// 重みづけデータの取得
		ticketDailyWeights := ticketDailyWeightRep.FindByTicketId(*ticket.Id)
		// 重みづけデータのticketDailyWeightValidMapの取得
		ticketDailyWeightValidMap := lo.Associate(ticketDailyWeights, func(item db.TicketDailyWeight) (int32, db.TicketDailyWeight) {
			return getValidIndex(globalStartDate, item.Date), item
		})
		if len(ticket.UserIds) > 0 {
			return &WeightWorkPerDayWithAssignCalculator{ticketDailyWeightRep, validUserMap, ticketDailyWeightValidMap}
		} else {
			return &WeightWorkPerDayWithNoAssignCalculator{ticketDailyWeightRep, validUserMap, ticketDailyWeightValidMap}
		}
	} else {
		return &DefaultWorkPerDayCalculator{}
	}
}

// DefaultWorkPerDayCalculator 通常のWorkPerDay計算（１日当たりの稼働量）
type DefaultWorkPerDayCalculator struct {
}

func (d *DefaultWorkPerDayCalculator) CalculateWorkPerDay(ticket *db.DefaultPileUp, validUserMap map[int32][]int32) float32 {
	// 重みづけのためにDB取得からフロントと同様ロジックで算出するように変更
	numberOfWorkerByDay := lo.Reduce(ticket.ValidIndexes, func(agg int, item int32, index int) int {
		validIndexUsers, ok := validUserMap[item]
		if !ok {
			return agg
		}
		return agg + len(lo.Filter(validIndexUsers, func(item int32, index int) bool {
			return lo.Contains(ticket.UserIds, item)
		}))
	}, 0)

	// アサイン済みの場合は一日当たりの稼働は総工数 / 総稼動可能日数
	return float32(*ticket.Estimate) / float32(numberOfWorkerByDay)
}

func (d *DefaultWorkPerDayCalculator) ReCalculateWorkPerDayByValidIndex(workHour float32, validIndex int32, ticket *db.DefaultPileUp) float32 {
	// デフォルトは何もなし
	return workHour
}

// WeightWorkPerDayWithAssignCalculator 重みづけ付のWorkPerDay計算（１日当たりの稼働量）
type WeightWorkPerDayWithAssignCalculator struct {
	ticketDailyWeightRep      interfaces.TicketDailyWeightRepositoryIF
	validUserMap              map[int32][]int32
	ticketDailyWeightValidMap map[int32]db.TicketDailyWeight
}

func (d *WeightWorkPerDayWithAssignCalculator) CalculateWorkPerDay(ticket *db.DefaultPileUp, validUserMap map[int32][]int32) float32 {
	// 重みづけデータの取得
	ticketDailyWeights := d.ticketDailyWeightRep.FindByTicketId(*ticket.Id)
	// 重みづけの日を除いた、総稼働人数を計算しなおす
	numberOfWorkerByDay := lo.Reduce(ticket.ValidIndexes, func(agg int, item int32, index int) int {

		// 重みづけの日は除外する
		if _, ok := d.ticketDailyWeightValidMap[item]; ok {
			return agg
		}

		validIndexUsers, ok := validUserMap[item]
		if !ok {
			return agg
		}
		return agg + len(lo.Intersect(validIndexUsers, ticket.UserIds))
	}, 0)

	// 重みづけで消費する工数を計算
	totalWeight := lo.Reduce(ticketDailyWeights, func(agg int32, item db.TicketDailyWeight, index int) int32 {
		return agg + item.WorkHour
	}, 0)

	// (予定工数 - 重みづけ工数) / 重みづけの日を除いた総稼働日数 で１日当たりの工数を算出。（重みづけの日を覗いて通常の計算をするイメージ）
	return float32(*ticket.Estimate-totalWeight) / float32(numberOfWorkerByDay)
}

func (d *WeightWorkPerDayWithAssignCalculator) ReCalculateWorkPerDayByValidIndex(workHour float32, validIndex int32, ticket *db.DefaultPileUp) float32 {
	// 重みづけが存在する日の場合
	if ticketDailyWeight, ok := d.ticketDailyWeightValidMap[validIndex]; ok {
		// その日の有効なユーザー一覧
		if validUsers, ok := d.validUserMap[validIndex]; ok {
			// 重みをその日有効なチケット担当者で割り振る
			workHour = float32(ticketDailyWeight.WorkHour) / float32(len(lo.Intersect(validUsers, ticket.UserIds)))
		}
	}
	return workHour
}

// WeightWorkPerDayWithNoAssignCalculator 重みづけ付のWorkPerDay計算（１日当たりの稼働量）
type WeightWorkPerDayWithNoAssignCalculator struct {
	ticketDailyWeightRep      interfaces.TicketDailyWeightRepositoryIF
	validUserMap              map[int32][]int32
	ticketDailyWeightValidMap map[int32]db.TicketDailyWeight
}

func (d *WeightWorkPerDayWithNoAssignCalculator) CalculateWorkPerDay(ticket *db.DefaultPileUp, validUserMap map[int32][]int32) float32 {
	// 重みづけデータの取得
	ticketDailyWeights := d.ticketDailyWeightRep.FindByTicketId(*ticket.Id)
	// 重みづけデータの工数をすべて足す
	numberOfWorkerByDay := lo.Reduce(ticketDailyWeights, func(agg int32, item db.TicketDailyWeight, index int) int32 {
		return agg + item.WorkHour
	}, 0)
	// 重みづけの日を除いた、総稼働人数を計算しなおす
	numberOfWorkDay := len(lo.Filter(ticket.ValidIndexes, func(item int32, index int) bool {
		// 重みづけの日は除外する
		if _, ok := d.ticketDailyWeightValidMap[item]; ok {
			return false
		}
		return true
	}))
	// 予定工数から重みづけの総数を除外し分配する
	return float32(*ticket.Estimate-numberOfWorkerByDay) / float32(numberOfWorkDay) / float32(*ticket.NumberOfWorker)
}

func (d *WeightWorkPerDayWithNoAssignCalculator) ReCalculateWorkPerDayByValidIndex(workHour float32, validIndex int32, ticket *db.DefaultPileUp) float32 {
	// 重みづけが存在する日の場合
	if ticketDailyWeight, ok := d.ticketDailyWeightValidMap[validIndex]; ok {
		// 未アサインの場合は人のように分配が不要なので重みをそのまま利用する
		return float32(ticketDailyWeight.WorkHour)
	}
	return workHour
}
