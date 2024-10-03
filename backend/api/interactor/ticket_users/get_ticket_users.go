package ticket_users

import (
	"github.com/gin-gonic/gin"
	"github.com/kenkonno/gantt-chart-proto/backend/api/middleware"
	"github.com/kenkonno/gantt-chart-proto/backend/api/openapi_models"
	"github.com/kenkonno/gantt-chart-proto/backend/models/db"
	"github.com/kenkonno/gantt-chart-proto/backend/repository"
	"github.com/samber/lo"
	"strconv"
)

func GetTicketUsersInvoke(c *gin.Context) openapi_models.GetTicketUsersResponse {
	ticketUserRep := repository.NewTicketUserRepository(middleware.GetRepositoryMode(c)...)

	ticketIdParam := c.QueryArray("ticketIds")
	var ticketIds []int32
	for _, v := range ticketIdParam {
		vv, err := strconv.Atoi(v)
		if err != nil {
			panic(err)
		}
		ticketIds = append(ticketIds, int32(vv))
	}

	ticketUserList := ticketUserRep.FindAllByTicketIds(ticketIds)

	return openapi_models.GetTicketUsersResponse{
		List: lo.Map(ticketUserList, func(item db.TicketUser, index int) openapi_models.TicketUser {
			return openapi_models.TicketUser{
				Id:        item.Id,
				TicketId:  item.TicketId,
				UserId:    item.UserId,
				Order:     int32(item.Order),
				CreatedAt: item.CreatedAt,
				UpdatedAt: item.UpdatedAt,
			}
		}),
	}
}
