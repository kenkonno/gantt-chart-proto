package ticket_daily_weights

import (
	"github.com/gin-gonic/gin"
	"github.com/kenkonno/gantt-chart-proto/backend/api/middleware"
	"github.com/kenkonno/gantt-chart-proto/backend/api/openapi_models"
	"github.com/kenkonno/gantt-chart-proto/backend/repository"
)

func DeleteTicketDailyWeightsIdInvoke(c *gin.Context) (openapi_models.DeleteTicketDailyWeightsIdResponse, error) {

	ticketDailyWeightRep := repository.NewTicketDailyWeightRepository(middleware.GetRepositoryMode(c)...)

	var req openapi_models.DeleteTicketDailyWeightsIdRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		return openapi_models.DeleteTicketDailyWeightsIdResponse{}, err
	}

	ticketDailyWeightRep.Delete(req.TicketId, req.Date)

	return openapi_models.DeleteTicketDailyWeightsIdResponse{}, nil
}
