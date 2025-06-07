package ticket_daily_weights

import (
	"github.com/gin-gonic/gin"
	"github.com/kenkonno/gantt-chart-proto/backend/api/middleware"
	"github.com/kenkonno/gantt-chart-proto/backend/api/openapi_models"
	"github.com/kenkonno/gantt-chart-proto/backend/repository"
	"strconv"
)

func GetTicketDailyWeightsIdInvoke(c *gin.Context) (openapi_models.GetTicketDailyWeightsIdResponse, error) {
	ticketDailyWeightRep := repository.NewTicketDailyWeightRepository(middleware.GetRepositoryMode(c)...)

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		panic(err)
	}

	ticketDailyWeight := ticketDailyWeightRep.Find(int32(id))

	return openapi_models.GetTicketDailyWeightsIdResponse{
		TicketDailyWeight: openapi_models.TicketDailyWeight{
			Id:        ticketDailyWeight.Id,
			TicketId:  ticketDailyWeight.TicketId,
			WorkHour:  ticketDailyWeight.WorkHour,
			Date:      ticketDailyWeight.Date,
			CreatedAt: ticketDailyWeight.CreatedAt,
			UpdatedAt: ticketDailyWeight.UpdatedAt,
		},
	}, nil
}
