package ticket_daily_weights

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kenkonno/gantt-chart-proto/backend/api/middleware"
	"github.com/kenkonno/gantt-chart-proto/backend/api/openapi_models"
	"github.com/kenkonno/gantt-chart-proto/backend/models/db"
	"github.com/kenkonno/gantt-chart-proto/backend/repository"
	"github.com/samber/lo"
)

func GetTicketDailyWeightsInvoke(c *gin.Context) (openapi_models.GetTicketDailyWeightsResponse, error) {
	ticketDailyWeightRep := repository.NewTicketDailyWeightRepository(middleware.GetRepositoryMode(c)...)

	var ticketDailyWeightList []db.TicketDailyWeight

	// facilityIdは必須項目ではない
	facilityIdStr := c.Query("facilityId")
	if facilityIdStr != "" {
		// facilityIdが存在する場合はFindByFacilityIdで検索
		facilityId, err := strconv.Atoi(facilityIdStr)
		if err != nil {
			panic(err)
		}
		ticketDailyWeightList = ticketDailyWeightRep.FindByFacilityId(int32(facilityId))
	} else {
		// facilityIdが存在しない場合はFindAllで検索
		ticketDailyWeightList = ticketDailyWeightRep.FindAll()
	}

	return openapi_models.GetTicketDailyWeightsResponse{
		List: lo.Map(ticketDailyWeightList, func(item db.TicketDailyWeight, index int) openapi_models.TicketDailyWeight {
			return openapi_models.TicketDailyWeight{
				TicketId:  item.TicketId,
				WorkHour:  &item.WorkHour,
				Date:      item.Date,
				CreatedAt: item.CreatedAt,
				UpdatedAt: item.UpdatedAt,
			}
		}),
	}, nil
}
