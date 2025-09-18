package ticket_daily_weights

import (
	"github.com/gin-gonic/gin"
	"github.com/kenkonno/gantt-chart-proto/backend/api/middleware"
	"github.com/kenkonno/gantt-chart-proto/backend/api/openapi_models"
	"github.com/kenkonno/gantt-chart-proto/backend/models/db"
	"github.com/kenkonno/gantt-chart-proto/backend/repository"
	"net/http"
	"time"
)

func PostTicketDailyWeightsInvoke(c *gin.Context) (openapi_models.PostTicketDailyWeightsResponse, error) {

	ticketDailyWeightRep := repository.NewTicketDailyWeightRepository(middleware.GetRepositoryMode(c)...)

	var ticketDailyWeightReq openapi_models.PostTicketDailyWeightsRequest
	if err := c.ShouldBindJSON(&ticketDailyWeightReq); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		panic(err)
	}

	// WorkHourがnilの時はそもそもレコードが不要なので削除する
	if ticketDailyWeightReq.TicketDailyWeight.WorkHour != nil {
		ticketDailyWeightRep.Upsert(db.TicketDailyWeight{
			TicketId:  ticketDailyWeightReq.TicketDailyWeight.TicketId,
			WorkHour:  *ticketDailyWeightReq.TicketDailyWeight.WorkHour,
			Date:      ticketDailyWeightReq.TicketDailyWeight.Date,
			CreatedAt: time.Time{},
			UpdatedAt: 0,
		})
	} else {
		ticketDailyWeightRep.Delete(ticketDailyWeightReq.TicketDailyWeight.TicketId, ticketDailyWeightReq.TicketDailyWeight.Date)
	}

	return openapi_models.PostTicketDailyWeightsResponse{}, nil
}
