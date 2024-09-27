package ticket_memo

import (
	"github.com/gin-gonic/gin"
	"github.com/kenkonno/gantt-chart-proto/backend/api/middleware"
	"github.com/kenkonno/gantt-chart-proto/backend/api/openapi_models"
	"github.com/kenkonno/gantt-chart-proto/backend/repository"
	"strconv"
)

func GetTicketMemoIdInvoke(c *gin.Context) openapi_models.GetTicketMemoIdResponse {
	ticketRep := repository.NewTicketRepository(middleware.GetRepositoryMode(c)...)

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		panic(err)
	}

	ticket := ticketRep.Find(int32(id))

	return openapi_models.GetTicketMemoIdResponse{
		Memo: ticket.Memo,
	}
}
