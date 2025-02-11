package milestones

import (
	"github.com/gin-gonic/gin"
	"github.com/kenkonno/gantt-chart-proto/backend/api/middleware"
	"github.com/kenkonno/gantt-chart-proto/backend/api/openapi_models"
	"github.com/kenkonno/gantt-chart-proto/backend/repository"
	"strconv"
)

func GetMilestonesIdInvoke(c *gin.Context) (openapi_models.GetMilestonesIdResponse, error) {
	milestoneRep := repository.NewMilestoneRepository(middleware.GetRepositoryMode(c)...)

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		panic(err)
	}

	milestone := milestoneRep.Find(int32(id))

	return openapi_models.GetMilestonesIdResponse{
		Milestone: openapi_models.Milestone{
			Id:              milestone.Id,
			FacilityId:      milestone.FacilityId,
			Date:            milestone.Date,
			Description:     milestone.Description,
			Order:           int32(milestone.Order),
			CreatedAt:       milestone.CreatedAt,
			UpdatedAt:       int(milestone.UpdatedAt),
		},
	}, nil
}
