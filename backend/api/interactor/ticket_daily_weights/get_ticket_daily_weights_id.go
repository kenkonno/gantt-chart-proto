package ticket_daily_weights

import (
	"github.com/gin-gonic/gin"
	"github.com/kenkonno/gantt-chart-proto/backend/api/middleware"
	"github.com/kenkonno/gantt-chart-proto/backend/api/openapi_models"
	"github.com/kenkonno/gantt-chart-proto/backend/repository"
)

func GetTicketDailyWeightsIdInvoke(c *gin.Context) (openapi_models.GetTicketDailyWeightsIdResponse, error) {
	ticketDailyWeightRep := repository.NewTicketDailyWeightRepository(middleware.GetRepositoryMode(c)...)

	var req openapi_models.GetTicketDailyWeightsIdRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		return openapi_models.GetTicketDailyWeightsIdResponse{}, err
	}

	ticketDailyWeight := ticketDailyWeightRep.Find(req.TicketId, req.Date)

	return openapi_models.GetTicketDailyWeightsIdResponse{
		TicketDailyWeight: openapi_models.TicketDailyWeight{
			TicketId:  ticketDailyWeight.TicketId,
			WorkHour: &ticketDailyWeight.WorkHour,
			Date:      ticketDailyWeight.Date,
			CreatedAt: ticketDailyWeight.CreatedAt,
			UpdatedAt: ticketDailyWeight.UpdatedAt,
		},
	}, nil
}
