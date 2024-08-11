package ticket_memo

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/kenkonno/gantt-chart-proto/backend/api/openapi_models"
	"github.com/kenkonno/gantt-chart-proto/backend/models/db"
	"github.com/kenkonno/gantt-chart-proto/backend/repository"
	"github.com/kenkonno/gantt-chart-proto/backend/repository/connection"
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

	result, err := ticketRep.UpdateMemo(db.Ticket{
		Id:   &id32,
		Memo: ticketReq.Memo,
		UpdatedAt: ticketReq.UpdatedAt,
	})
	if err != nil {
		var target connection.ConflictError
		if errors.As(err, &target) {
			c.JSON(http.StatusConflict, err.Error())
			panic(err)
		}
	}
	return openapi_models.PostTicketMemoIdResponse{
		Msg:       result.Memo,
		UpdatedAt: result.UpdatedAt,
	}

}
