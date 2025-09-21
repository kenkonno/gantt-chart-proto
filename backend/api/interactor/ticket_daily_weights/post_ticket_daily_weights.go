package ticket_daily_weights

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kenkonno/gantt-chart-proto/backend/api/middleware"
	"github.com/kenkonno/gantt-chart-proto/backend/api/openapi_models"
	"github.com/kenkonno/gantt-chart-proto/backend/models/db"
	"github.com/kenkonno/gantt-chart-proto/backend/repository"
	"github.com/samber/lo"
)

func PostTicketDailyWeightsInvoke(c *gin.Context) (openapi_models.PostTicketDailyWeightsResponse, error) {

	ticketDailyWeightRep := repository.NewTicketDailyWeightRepository(middleware.GetRepositoryMode(c)...)
	ticketRep := repository.NewTicketRepository(middleware.GetRepositoryMode(c)...)

	var ticketDailyWeightReq openapi_models.PostTicketDailyWeightsRequest
	if err := c.ShouldBindJSON(&ticketDailyWeightReq); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		panic(err)
	}

	if ticketDailyWeightReq.TicketDailyWeight.WorkHour != nil {
		ticketDailyWeightRep.Upsert(db.TicketDailyWeight{
			TicketId:  ticketDailyWeightReq.TicketDailyWeight.TicketId,
			WorkHour:  *ticketDailyWeightReq.TicketDailyWeight.WorkHour,
			Date:      ticketDailyWeightReq.TicketDailyWeight.Date,
			CreatedAt: time.Time{},
			UpdatedAt: 0,
		})
		// チケットに紐づく重みづけをすべて取得
		dailyWeights := ticketDailyWeightRep.FindByTicketId(ticketDailyWeightReq.TicketDailyWeight.TicketId)
		totalWeights := lo.SumBy(dailyWeights, func(item db.TicketDailyWeight) int32 {
			return item.WorkHour
		})
		estimate := ticketRep.Find(ticketDailyWeightReq.TicketDailyWeight.TicketId).Estimate

		if estimate != nil && totalWeights > *estimate {
			// TODO: 自動生成の対応が追いついていないので修正する
			err := errors.New("総工数を上回る重みづけが入力されています。")
			c.JSON(http.StatusBadRequest, "総工数を上回る重みづけが入力されています。")
			return openapi_models.PostTicketDailyWeightsResponse{}, err
		}

	} else {
		// WorkHourがnilの時はそもそもレコードが不要なので削除する
		ticketDailyWeightRep.Delete(ticketDailyWeightReq.TicketDailyWeight.TicketId, ticketDailyWeightReq.TicketDailyWeight.Date)
	}

	return openapi_models.PostTicketDailyWeightsResponse{}, nil
}
