package milestones

import (
	"github.com/gin-gonic/gin"
	"github.com/kenkonno/gantt-chart-proto/backend/api/openapi_models"
	"github.com/kenkonno/gantt-chart-proto/backend/models/db"
	"github.com/kenkonno/gantt-chart-proto/backend/repository"
	"net/http"
	"strings"
	"time"
)

func PostMilestonesIdInvoke(c *gin.Context) openapi_models.PostMilestonesIdResponse {

	milestoneRep := repository.NewMilestoneRepository()

	var milestoneReq openapi_models.PostMilestonesRequest
	if err := c.ShouldBindJSON(&milestoneReq); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		panic(err)
	}

	milestoneRep.Upsert(db.Milestone{
		Id:              milestoneReq.Milestone.Id,
		FacilityId:      milestoneReq.Milestone.FacilityId,
		Date:            milestoneReq.Milestone.Date,
		Description:     strings.TrimSpace(milestoneReq.Milestone.Description),
		Order:           int(milestoneReq.Milestone.Order),
		CreatedAt:       time.Time{},
		UpdatedAt:       0,
	})

	return openapi_models.PostMilestonesIdResponse{}

}
