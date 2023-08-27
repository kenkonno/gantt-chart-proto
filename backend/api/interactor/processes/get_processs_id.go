package processes

import (
	"github.com/gin-gonic/gin"
	"github.com/kenkonno/gantt-chart-proto/backend/api/openapi_models"
	"github.com/kenkonno/gantt-chart-proto/backend/repository"
	"strconv"
)

func GetProcessesIdInvoke(c *gin.Context) openapi_models.GetProcessesIdResponse {
	processRep := repository.NewProcessRepository()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		panic(err)
	}

	process := processRep.Find(int32(id))

	return openapi_models.GetProcessesIdResponse{
		Process: openapi_models.Process{
			Id:        process.Id,
			Name:      process.Name,
			CreatedAt: process.CreatedAt,
			UpdatedAt: process.UpdatedAt,
		},
	}
}
