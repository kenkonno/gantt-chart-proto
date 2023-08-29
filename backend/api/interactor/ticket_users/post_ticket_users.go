package ticket_users

import (
	"github.com/gin-gonic/gin"
	"github.com/kenkonno/gantt-chart-proto/backend/api/openapi_models"
	"github.com/kenkonno/gantt-chart-proto/backend/models/db"
	"github.com/kenkonno/gantt-chart-proto/backend/repository"
	"net/http"
	"time"
)

// PostTicketUsersInvoke このAPIはticketIdを受け取ってユーザーを更新するものとする。
func PostTicketUsersInvoke(c *gin.Context) openapi_models.PostTicketUsersResponse {

	ticketUserRep := repository.NewTicketUserRepository()
	var ticketUserReq openapi_models.PostTicketUsersRequest
	if err := c.ShouldBindJSON(&ticketUserReq); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		panic(err)
	}

	var result []openapi_models.TicketUser
	ticketUserRep.DeleteByTicketId(ticketUserReq.TicketId)
	for _, v := range ticketUserReq.UserIds {
		r := ticketUserRep.Upsert(db.TicketUser{
			TicketId:  ticketUserReq.TicketId,
			UserId:    v,
			CreatedAt: time.Time{},
			UpdatedAt: 0,
		})
		result = append(result, openapi_models.TicketUser{
			Id:        r.Id,
			TicketId:  r.TicketId,
			UserId:    r.UserId,
			CreatedAt: r.CreatedAt,
			UpdatedAt: r.UpdatedAt,
		})
	}

	return openapi_models.PostTicketUsersResponse{
		TicketUsers: result,
	}

}
