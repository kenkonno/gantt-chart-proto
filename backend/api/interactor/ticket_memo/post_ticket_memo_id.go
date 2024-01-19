package ticket_memo

import (
	"github.com/gin-gonic/gin"
	"github.com/kenkonno/gantt-chart-proto/backend/api/openapi_models"
	"github.com/kenkonno/gantt-chart-proto/backend/models/db"
	"github.com/kenkonno/gantt-chart-proto/backend/repository"
	"net/http"
	"strconv"
)

func PostTicketMemoIdInvoke(c *gin.Context) openapi_models.PostTicketMemoIdResponse {

	ticketRep := repository.NewTicketRepository()

	id, err := strconv.Atoi(c.Param("id"))
	id32 := int32(id)
	if err != nil {
		panic(err)
	}

	var ticketReq openapi_models.PostTicketMemoIdRequest
	if err := c.ShouldBindJSON(&ticketReq); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		panic(err)
	}

	ticketRep.UpdateMemo(db.Ticket{
		Id:   &id32,
		Memo: ticketReq.Memo,
	})
	return openapi_models.PostTicketMemoIdResponse{}

}
