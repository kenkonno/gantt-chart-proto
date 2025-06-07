package ticket_daily_weights

import (
	"github.com/gin-gonic/gin"
	"github.com/kenkonno/gantt-chart-proto/backend/api/middleware"
	"github.com/kenkonno/gantt-chart-proto/backend/api/openapi_models"
	"github.com/kenkonno/gantt-chart-proto/backend/repository"
	"strconv"
)

func DeleteTicketDailyWeightsIdInvoke(c *gin.Context) (openapi_models.DeleteTicketDailyWeightsIdResponse, error) {

	ticketDailyWeightRep := repository.NewTicketDailyWeightRepository(middleware.GetRepositoryMode(c)...)

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		panic(err)
	}

	ticketDailyWeightRep.Delete(int32(id))

	return openapi_models.DeleteTicketDailyWeightsIdResponse{}, nil
}
