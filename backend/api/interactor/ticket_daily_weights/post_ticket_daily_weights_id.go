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

// 未使用
func PostTicketDailyWeightsIdInvoke(c *gin.Context) (openapi_models.PostTicketDailyWeightsIdResponse, error) {

	ticketDailyWeightRep := repository.NewTicketDailyWeightRepository(middleware.GetRepositoryMode(c)...)

	var ticketDailyWeightReq openapi_models.PostTicketDailyWeightsRequest
	if err := c.ShouldBindJSON(&ticketDailyWeightReq); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		panic(err)
	}

	ticketDailyWeightRep.Upsert(db.TicketDailyWeight{
		TicketId:  ticketDailyWeightReq.TicketDailyWeight.TicketId,
		WorkHour: *ticketDailyWeightReq.TicketDailyWeight.WorkHour,
		Date:      ticketDailyWeightReq.TicketDailyWeight.Date,
		CreatedAt: time.Time{},
		UpdatedAt: 0,
	})

	return openapi_models.PostTicketDailyWeightsIdResponse{}, nil
}
