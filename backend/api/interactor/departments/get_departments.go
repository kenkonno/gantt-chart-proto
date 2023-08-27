package departments

import (
	"github.com/gin-gonic/gin"
	"github.com/kenkonno/gantt-chart-proto/backend/api/openapi_models"
	"github.com/kenkonno/gantt-chart-proto/backend/models/db"
	"github.com/kenkonno/gantt-chart-proto/backend/repository"
	"github.com/samber/lo"
)

func GetDepartmentsInvoke(c *gin.Context) openapi_models.GetDepartmentsResponse {
	departmentRep := repository.NewDepartmentRepository()

	departmentList := departmentRep.FindAll()

	return openapi_models.GetDepartmentsResponse{
		List: lo.Map(departmentList, func(item db.Department, index int) openapi_models.Department {
			return openapi_models.Department{
				Id:        item.Id,
				Name:      item.Name,
				CreatedAt: item.CreatedAt,
				UpdatedAt: item.UpdatedAt,
			}
		}),
	}
}
