package holidays

import (
	"github.com/gin-gonic/gin"
	"github.com/kenkonno/gantt-chart-proto/backend/api/middleware"
	"github.com/kenkonno/gantt-chart-proto/backend/api/openapi_models"
	"github.com/kenkonno/gantt-chart-proto/backend/models/db"
	"github.com/kenkonno/gantt-chart-proto/backend/repository"
	"github.com/samber/lo"
)

func GetHolidaysInvoke(c *gin.Context) (openapi_models.GetHolidaysResponse, error) {
	holidayRep := repository.NewHolidayRepository(middleware.GetRepositoryMode(c)...)

	holidayList := holidayRep.FindAll()

	return openapi_models.GetHolidaysResponse{
		List: lo.Map(holidayList, func(item db.Holiday, index int) openapi_models.Holiday {
			return openapi_models.Holiday{
				Id:        item.Id,
				Name:      item.Name,
				Date:      item.Date,
				CreatedAt: item.CreatedAt,
				UpdatedAt: item.UpdatedAt,
			}
		}),
	}, nil
}
