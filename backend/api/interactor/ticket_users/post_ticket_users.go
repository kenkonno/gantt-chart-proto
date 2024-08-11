package ticket_users

import (
	"github.com/gin-gonic/gin"
	"github.com/kenkonno/gantt-chart-proto/backend/api/openapi_models"
	"github.com/kenkonno/gantt-chart-proto/backend/models/db"
	"github.com/kenkonno/gantt-chart-proto/backend/repository"
	"github.com/kenkonno/gantt-chart-proto/backend/repository/connection"
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

	// DeleteInsertなのでInteractor側で CreatedAtが同一であることを確認する。
	originals := ticketUserRep.FindByTicketId(ticketUserReq.TicketId)
	if len(originals) > 0 {
		if originals[0].CreatedAt.Truncate(time.Second) != ticketUserReq.CreatedAt.Truncate(time.Second) {
			err := connection.NewConflictError()
			c.JSON(http.StatusConflict, err.Error())
			panic(err)
		}
	}

	var result = []openapi_models.TicketUser{}
	ticketUserRep.DeleteByTicketId(ticketUserReq.TicketId)
	createdAt := time.Now()
	for index, v := range ticketUserReq.UserIds {
		r := ticketUserRep.UpsertWithCreatedAt(db.TicketUser{
			TicketId:  ticketUserReq.TicketId,
			UserId:    v,
			Order:     index,
			CreatedAt: createdAt,
			UpdatedAt: 0,
		})
		result = append(result, openapi_models.TicketUser{
			Id:        r.Id,
			TicketId:  r.TicketId,
			UserId:    r.UserId,
			Order:     int32(r.Order),
			CreatedAt: r.CreatedAt,
			UpdatedAt: r.UpdatedAt,
		})
	}

	return openapi_models.PostTicketUsersResponse{
		TicketUsers: result,
	}

}
