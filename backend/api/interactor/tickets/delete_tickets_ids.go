package tickets

import (
	"github.com/gin-gonic/gin"
	"github.com/kenkonno/gantt-chart-proto/backend/api/middleware"
	"github.com/kenkonno/gantt-chart-proto/backend/api/openapi_models"
	"github.com/kenkonno/gantt-chart-proto/backend/repository"
	"strconv"
)

func DeleteTicketsIdInvoke(c *gin.Context) (openapi_models.DeleteTicketsIdResponse, error) {

	ticketRep := repository.NewTicketRepository(middleware.GetRepositoryMode(c)...)

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		panic(err)
	}

	ticketRep.Delete(int32(id))

	return openapi_models.DeleteTicketsIdResponse{}, nil

}
