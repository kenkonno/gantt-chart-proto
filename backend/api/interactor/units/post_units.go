package units

import (
	"github.com/gin-gonic/gin"
	"github.com/kenkonno/gantt-chart-proto/backend/api/openapi_models"
	"github.com/kenkonno/gantt-chart-proto/backend/models/db"
	"github.com/kenkonno/gantt-chart-proto/backend/repository"
	"net/http"
	"time"
)

func PostUnitsInvoke(c *gin.Context) openapi_models.PostUnitsResponse {

	unitRep := repository.NewUnitRepository()
	ganttGroupsRep := repository.NewGanttGroupRepository()

	var unitReq openapi_models.PostUnitsRequest
	if err := c.ShouldBindJSON(&unitReq); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		panic(err)
	}
	r := unitRep.Upsert(db.Unit{
		Name:       unitReq.Unit.Name,
		FacilityId: unitReq.Unit.FacilityId,
		CreatedAt:  time.Time{},
		UpdatedAt:  0,
	})

	ganttGroupsRep.Upsert(db.GanttGroup{
		Id:         nil,
		FacilityId: unitReq.Unit.FacilityId,
		UnitId:     *r.Id,
		Order:      0, // TODO: ユニット一覧に並び順を持たせるべき 今はもう適当にする
		CreatedAt:  time.Time{},
		UpdatedAt:  0,
	})

	return openapi_models.PostUnitsResponse{}

}
