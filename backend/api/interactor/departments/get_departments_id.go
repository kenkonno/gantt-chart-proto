package departments

import (
	"github.com/gin-gonic/gin"
	"github.com/kenkonno/gantt-chart-proto/backend/api/middleware"
	"github.com/kenkonno/gantt-chart-proto/backend/api/openapi_models"
	"github.com/kenkonno/gantt-chart-proto/backend/repository"
	"strconv"
)

func GetDepartmentsIdInvoke(c *gin.Context) (openapi_models.GetDepartmentsIdResponse, error) {
	departmentRep := repository.NewDepartmentRepository(middleware.GetRepositoryMode(c)...)

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		panic(err)
	}

	department := departmentRep.Find(int32(id))

	return openapi_models.GetDepartmentsIdResponse{
		Department: openapi_models.Department{
			Id:        department.Id,
			Name:      department.Name,
			Color:     department.Color,
			Order:     int32(department.Order),
			CreatedAt: department.CreatedAt,
			UpdatedAt: department.UpdatedAt,
		},
	}, nil
}
