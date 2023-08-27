package departments

import (
	"github.com/gin-gonic/gin"
	"github.com/kenkonno/gantt-chart-proto/backend/api/openapi_models"
	"github.com/kenkonno/gantt-chart-proto/backend/repository"
	"strconv"
)

func DeleteDepartmentsIdInvoke(c *gin.Context) openapi_models.DeleteDepartmentsIdResponse {

	departmentRep := repository.NewDepartmentRepository()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		panic(err)
	}

	departmentRep.Delete(int32(id))

	return openapi_models.DeleteDepartmentsIdResponse{}

}
