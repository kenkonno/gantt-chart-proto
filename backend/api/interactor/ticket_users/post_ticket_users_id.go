package ticket_users

import (
	"github.com/gin-gonic/gin"
	"github.com/kenkonno/gantt-chart-proto/backend/api/openapi_models"
	"github.com/kenkonno/gantt-chart-proto/backend/models/db"
	"github.com/kenkonno/gantt-chart-proto/backend/repository"
	"net/http"
	"time"
)

func PostTicketUsersIdInvoke(c *gin.Context) openapi_models.PostTicketUsersIdResponse {

	ticketUserRep := repository.NewTicketUserRepository()

	var ticketUserReq openapi_models.PostTicketUsersRequest
	if err := c.ShouldBindJSON(&ticketUserReq); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		panic(err)
	}

	ticketUserRep.Upsert(db.TicketUser{
		Id:        ticketUserReq.TicketUser.Id,
		TicketId:  ticketUserReq.TicketUser.TicketId,
		UserId:    ticketUserReq.TicketUser.UserId,
		CreatedAt: time.Time{},
		UpdatedAt: 0,
	})

	return openapi_models.PostTicketUsersIdResponse{}

}
