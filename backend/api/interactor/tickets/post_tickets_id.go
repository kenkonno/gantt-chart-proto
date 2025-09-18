package tickets

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kenkonno/gantt-chart-proto/backend/api/constants"
	"github.com/kenkonno/gantt-chart-proto/backend/api/middleware"
	"github.com/kenkonno/gantt-chart-proto/backend/api/openapi_models"
	"github.com/kenkonno/gantt-chart-proto/backend/api/utils"
	"github.com/kenkonno/gantt-chart-proto/backend/models/db"
	"github.com/kenkonno/gantt-chart-proto/backend/repository"
	"github.com/kenkonno/gantt-chart-proto/backend/repository/connection"
	"github.com/kenkonno/gantt-chart-proto/backend/repository/interfaces"
)

func PostTicketsIdInvoke(c *gin.Context) (openapi_models.PostTicketsIdResponse, error) {

	ticketRep := repository.NewTicketRepository(middleware.GetRepositoryMode(c)...)
	ticketDailyWeightsRep := repository.NewTicketDailyWeightRepository(middleware.GetRepositoryMode(c)...)

	var ticketReq openapi_models.PostTicketsRequest
	if err := c.ShouldBindJSON(&ticketReq); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		panic(err)
	}

	originalTicket := ticketRep.Find(*ticketReq.Ticket.Id)

	result, err := ticketRep.Upsert(db.Ticket{
		Id:              ticketReq.Ticket.Id,
		GanttGroupId:    ticketReq.Ticket.GanttGroupId,
		ProcessId:       ticketReq.Ticket.ProcessId,
		DepartmentId:    ticketReq.Ticket.DepartmentId,
		LimitDate:       ticketReq.Ticket.LimitDate,
		Estimate:        ticketReq.Ticket.Estimate,
		NumberOfWorker:  ticketReq.Ticket.NumberOfWorker,
		DaysAfter:       ticketReq.Ticket.DaysAfter,
		StartDate:       ticketReq.Ticket.StartDate,
		EndDate:         ticketReq.Ticket.EndDate,
		ProgressPercent: ticketReq.Ticket.ProgressPercent,
		Order:           ticketReq.Ticket.Order,
		CreatedAt:       time.Time{},
		UpdatedAt:       ticketReq.Ticket.UpdatedAt,
	})
	if err != nil {
		var target connection.ConflictError
		if errors.As(err, &target) {
			c.JSON(http.StatusConflict, err.Error())
			panic(err)
		}
	}

	if utils.HasOption(constants.WorkloadWeighting) {
		updateTicketDailyWeights(ticketDailyWeightsRep, result, originalTicket)
	}

	return openapi_models.PostTicketsIdResponse{
		Ticket: openapi_models.Ticket{
			Id:              result.Id,
			GanttGroupId:    result.GanttGroupId,
			ProcessId:       result.ProcessId,
			DepartmentId:    result.DepartmentId,
			LimitDate:       result.LimitDate,
			Estimate:        result.Estimate,
			NumberOfWorker:  result.NumberOfWorker,
			DaysAfter:       result.DaysAfter,
			StartDate:       result.StartDate,
			EndDate:         result.EndDate,
			ProgressPercent: result.ProgressPercent,
			Memo:            result.Memo,
			Order:           result.Order,
			CreatedAt:       result.CreatedAt,
			UpdatedAt:       result.UpdatedAt,
		},
	}, nil

}

func updateTicketDailyWeights(ticketDailyWeightsRep interfaces.TicketDailyWeightRepositoryIF, result db.Ticket, originalTicket db.Ticket) {
	// TODO: バグってるポイし祝日の挙動が考慮されていないのでいったんコメントアウト
	return
	/* 重みづけデータがある場合は以下の仕様に基づき日付を移動させる。
	1. 期間の移動であれば以下の処理を行う
	  1. 更新前のチケットの開始日と重みづけデータの各日付を比較し、何日後なのかを計算する。
	  2. 重みづけデータを削除する。
	  3. 重みづけデータすべてを更新後のチケットの開始日からN日後に更新する。
	2. それ以外の場合は、チケットの期間外の重みづけデータのみを削除する
	*/
	// 重みづけデータを取得
	ticketDailyWeights := ticketDailyWeightsRep.FindByTicketId(*result.Id)

	// 重みづけデータが存在する場合のみ処理
	if len(ticketDailyWeights) > 0 && originalTicket.StartDate != nil && result.StartDate != nil {

		// 期間の移動かどうかを判定（更新前後で期間の日数が同じかどうか）
		var isPeriodMove bool
		if originalTicket.EndDate != nil && result.EndDate != nil {
			originalDays := int(originalTicket.EndDate.Sub(*originalTicket.StartDate).Hours()/24) + 1 // +1を追加（期間には開始日も含まれるため）
			newDays := int(result.EndDate.Sub(*result.StartDate).Hours()/24) + 1                      // +1を追加
			isPeriodMove = (originalDays == newDays)
		} else {
			isPeriodMove = false
		}

		if isPeriodMove {
			// 1. 期間の移動の場合
			for _, weight := range ticketDailyWeights {
				if !weight.Date.IsZero() {
					// 1-1. 更新前の開始日からの日数差を計算
					daysAfter := int(weight.Date.Sub(*originalTicket.StartDate).Hours() / 24)

					// 1-2. 元の重みづけデータを削除
					ticketDailyWeightsRep.Delete(*result.Id, weight.Date)

					// 1-3. 新しい日付で重みづけデータを作成（更新後のチケットの開始日からN日後）
					newDate := result.StartDate.AddDate(0, 0, daysAfter)
					newWeight := db.TicketDailyWeight{
						TicketId:  weight.TicketId,
						Date:      newDate,
						WorkHour:  weight.WorkHour,
						CreatedAt: time.Now(),
						UpdatedAt: int32(time.Now().Unix()),
					}
					ticketDailyWeightsRep.Upsert(newWeight)
				}
			}
		} else {
			// 2. それ以外の場合は、チケットの期間外の重みづけデータのみを削除
			for _, weight := range ticketDailyWeights {
				if !weight.Date.IsZero() {
					weightDate := weight.Date

					// チケットの期間外かどうかをチェック
					isOutsidePeriod := false
					if weightDate.Before(*result.StartDate) {
						isOutsidePeriod = true
					}
					if result.EndDate != nil && (weightDate.After(*result.EndDate) || weightDate.Equal(*result.EndDate)) {
						isOutsidePeriod = true
					}

					// 期間外の場合のみ削除
					if isOutsidePeriod {
						ticketDailyWeightsRep.Delete(*result.Id, weight.Date)
					}
				}
			}
		}
	}
}
