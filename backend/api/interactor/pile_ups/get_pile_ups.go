package pile_ups

import (
	"github.com/gin-gonic/gin"
	"github.com/kenkonno/gantt-chart-proto/backend/api/openapi_models"
	"github.com/kenkonno/gantt-chart-proto/backend/models/db"
	"github.com/kenkonno/gantt-chart-proto/backend/repository"
	"github.com/samber/lo"
	"strconv"
)

func GetPileUpsInvoke(c *gin.Context) openapi_models.GetPileUpsResponse {
	excludeFacilityId, err := strconv.Atoi(c.Query("facilityId"))
	if err != nil {
		panic(err)
	}
	facilityRep := repository.NewFacilityRepository()
	facilities := lo.Filter(facilityRep.FindAll(), func(item db.Facility, index int) bool {
		return *item.Id != int32(excludeFacilityId)
	})
	ganttGroupRep := repository.NewGanttGroupRepository()
	ganttGroups := ganttGroupRep.FindAll()
	holidayRep := repository.NewHolidayRepository()
	holidays := holidayRep.FindAll()

	return openapi_models.GetPileUpsResponse{
		List: lo.Map(facilities, func(facility db.Facility, index int) openapi_models.GetPileUpsResponseListInner {
			targetHolidays := lo.Filter(holidays, func(item db.Holiday, index int) bool {
				return item.FacilityId == *facility.Id
			})
			targetGanttGroups := lo.Filter(ganttGroups, func(item db.GanttGroup, index int) bool {
				return item.FacilityId == *facility.Id
			})
			return openapi_models.GetPileUpsResponseListInner{
				FacilityId: *facility.Id,
				Holidays: lo.Map(targetHolidays, func(item db.Holiday, index int) openapi_models.Holiday {
					return openapi_models.Holiday{
						Id:         item.Id,
						Name:       item.Name,
						Date:       item.Date,
						CreatedAt:  item.CreatedAt,
						UpdatedAt:  item.UpdatedAt,
						FacilityId: item.FacilityId,
					}
				}),
				GanttGroups: lo.Map(targetGanttGroups, func(item db.GanttGroup, index int) openapi_models.GanttGroup {
					return openapi_models.GanttGroup{
						Id:         item.Id,
						FacilityId: item.FacilityId,
						UnitId:     item.UnitId,
						CreatedAt:  item.CreatedAt,
						UpdatedAt:  item.UpdatedAt,
					}
				}),
			}
		}),
	}
}
