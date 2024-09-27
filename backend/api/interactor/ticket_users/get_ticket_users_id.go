package ticket_users

import (
	"github.com/gin-gonic/gin"
	"github.com/kenkonno/gantt-chart-proto/backend/api/middleware"
	"github.com/kenkonno/gantt-chart-proto/backend/api/openapi_models"
	"github.com/kenkonno/gantt-chart-proto/backend/repository"
	"strconv"
)

func GetTicketUsersIdInvoke(c *gin.Context) openapi_models.GetTicketUsersIdResponse {
	ticketUserRep := repository.NewTicketUserRepository(middleware.GetRepositoryMode(c)...)

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		panic(err)
	}

	ticketUser := ticketUserRep.Find(int32(id))

	return openapi_models.GetTicketUsersIdResponse{
		TicketUser: openapi_models.TicketUser{
			Id:        ticketUser.Id,
			TicketId:  ticketUser.TicketId,
			UserId:    ticketUser.UserId,
			Order:     int32(ticketUser.Order),
			CreatedAt: ticketUser.CreatedAt,
			UpdatedAt: ticketUser.UpdatedAt,
		},
	}
}
