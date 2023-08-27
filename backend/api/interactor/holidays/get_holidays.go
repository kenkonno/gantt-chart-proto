package holidays

import (
	"github.com/gin-gonic/gin"
	"github.com/kenkonno/gantt-chart-proto/backend/api/openapi_models"
	"github.com/kenkonno/gantt-chart-proto/backend/models/db"
	"github.com/kenkonno/gantt-chart-proto/backend/repository"
	"github.com/samber/lo"
	"strconv"
)

func GetHolidaysInvoke(c *gin.Context) openapi_models.GetHolidaysResponse {
	holidayRep := repository.NewHolidayRepository()
	facilityId, err := strconv.Atoi(c.Query("facilityId"))
	if err != nil {
		panic(err)
	}

	holidayList := holidayRep.FindByFacilityId(int32(facilityId))

	return openapi_models.GetHolidaysResponse{
		List: lo.Map(holidayList, func(item db.Holiday, index int) openapi_models.Holiday {
			return openapi_models.Holiday{
				Id:         item.Id,
				FacilityId: item.FacilityId,
				Name:       item.Name,
				Date:       item.Date,
				CreatedAt:  item.CreatedAt,
				UpdatedAt:  item.UpdatedAt,
			}
		}),
	}
}
