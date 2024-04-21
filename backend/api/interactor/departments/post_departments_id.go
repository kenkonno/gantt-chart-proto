package departments

import (
	"github.com/gin-gonic/gin"
	"github.com/kenkonno/gantt-chart-proto/backend/api/openapi_models"
	"github.com/kenkonno/gantt-chart-proto/backend/models/db"
	"github.com/kenkonno/gantt-chart-proto/backend/repository"
	"net/http"
	"strings"
	"time"
)

func PostDepartmentsIdInvoke(c *gin.Context) openapi_models.PostDepartmentsIdResponse {

	departmentRep := repository.NewDepartmentRepository()

	var departmentReq openapi_models.PostDepartmentsRequest
	if err := c.ShouldBindJSON(&departmentReq); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		panic(err)
	}

	departmentRep.Upsert(db.Department{
		Id:        departmentReq.Department.Id,
		Name:      strings.TrimSpace(departmentReq.Department.Name),
		Order:     int(departmentReq.Department.Order),
		CreatedAt: time.Time{},
		UpdatedAt: 0,
	})

	return openapi_models.PostDepartmentsIdResponse{}

}
