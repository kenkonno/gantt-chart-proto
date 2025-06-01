package milestones

import (
	"github.com/gin-gonic/gin"
	"github.com/kenkonno/gantt-chart-proto/backend/api/middleware"
	"github.com/kenkonno/gantt-chart-proto/backend/api/openapi_models"
	"github.com/kenkonno/gantt-chart-proto/backend/models/db"
	"github.com/kenkonno/gantt-chart-proto/backend/repository"
	"github.com/samber/lo"
	"strconv"
)

func GetMilestonesInvoke(c *gin.Context) (openapi_models.GetMilestonesResponse, error) {

	facilityIds := c.QueryArray("facilityIds")
	facilityIdsInt32 := lo.Map(facilityIds, func(item string, index int) int32 {
		v, _ := strconv.Atoi(item)
		return int32(v)
	})
	mode := c.Query("mode")
	milestoneRep := repository.NewMilestoneRepository(middleware.GetRepositoryMode(c)...)
	if mode == "prod" {
		milestoneRep = repository.NewMilestoneRepository()
	}

	milestoneList := milestoneRep.FindByFacilityId(facilityIdsInt32)

	return openapi_models.GetMilestonesResponse{
		List: lo.Map(milestoneList, func(item db.Milestone, index int) openapi_models.Milestone {
			return openapi_models.Milestone{
				Id:          item.Id,
				FacilityId:  item.FacilityId,
				Date:        item.Date,
				Description: item.Description,
				Order:       int32(item.Order),
				CreatedAt:   item.CreatedAt,
				UpdatedAt:   int(item.UpdatedAt),
			}
		}),
	}, nil
}
