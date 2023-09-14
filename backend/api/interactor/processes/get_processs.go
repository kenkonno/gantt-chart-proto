package processes

import (
	"github.com/gin-gonic/gin"
	"github.com/kenkonno/gantt-chart-proto/backend/api/openapi_models"
	"github.com/kenkonno/gantt-chart-proto/backend/models/db"
	"github.com/kenkonno/gantt-chart-proto/backend/repository"
	"github.com/samber/lo"
)

func GetProcessesInvoke(c *gin.Context) openapi_models.GetProcessesResponse {
	processRep := repository.NewProcessRepository()

	processList := processRep.FindAll()

	return openapi_models.GetProcessesResponse{
		List: lo.Map(processList, func(item db.Process, index int) openapi_models.Process {
			return openapi_models.Process{
				Id:        item.Id,
				Name:      item.Name,
				Order:     int32(item.Order),
				CreatedAt: item.CreatedAt,
				UpdatedAt: item.UpdatedAt,
			}
		}),
	}
}
