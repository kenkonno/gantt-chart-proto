package ticket_daily_weights

import (
	"github.com/gin-gonic/gin"
	"github.com/kenkonno/gantt-chart-proto/backend/api/middleware"
	"github.com/kenkonno/gantt-chart-proto/backend/api/openapi_models"
	"github.com/kenkonno/gantt-chart-proto/backend/models/db"
	"github.com/kenkonno/gantt-chart-proto/backend/repository"
	"github.com/samber/lo"
)

func GetTicketDailyWeightsInvoke(c *gin.Context) (openapi_models.GetTicketDailyWeightsResponse, error) {
	ticketDailyWeightRep := repository.NewTicketDailyWeightRepository(middleware.GetRepositoryMode(c)...)

	ticketDailyWeightList := ticketDailyWeightRep.FindAll()

	return openapi_models.GetTicketDailyWeightsResponse{
		List: lo.Map(ticketDailyWeightList, func(item db.TicketDailyWeight, index int) openapi_models.TicketDailyWeight {
			return openapi_models.TicketDailyWeight{
				Id:        item.Id,
				TicketId:  item.TicketId,
				WorkHour:  item.WorkHour,
				Date:      item.Date,
				CreatedAt: item.CreatedAt,
				UpdatedAt: item.UpdatedAt,
			}
		}),
	}, nil
}
