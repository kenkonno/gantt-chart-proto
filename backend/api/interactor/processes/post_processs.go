package processes

import (
	"github.com/gin-gonic/gin"
	"github.com/kenkonno/gantt-chart-proto/backend/api/openapi_models"
	"github.com/kenkonno/gantt-chart-proto/backend/models/db"
	"github.com/kenkonno/gantt-chart-proto/backend/repository"
	"net/http"
	"time"
)

func PostProcessesInvoke(c *gin.Context) openapi_models.PostProcessesResponse {

	processRep := repository.NewProcessRepository()

	var processReq openapi_models.PostProcessesRequest
	if err := c.ShouldBindJSON(&processReq); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		panic(err)
	}
	processRep.Upsert(db.Process{
		Name:      processReq.Process.Name,
		Order:     int(processReq.Process.Order),
		CreatedAt: time.Time{},
		UpdatedAt: 0,
	})

	return openapi_models.PostProcessesResponse{}

}
