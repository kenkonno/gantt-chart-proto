package processes

import (
	"github.com/gin-gonic/gin"
	"github.com/kenkonno/gantt-chart-proto/backend/api/openapi_models"
	"github.com/kenkonno/gantt-chart-proto/backend/repository"
	"strconv"
)

func DeleteProcessesIdInvoke(c *gin.Context) openapi_models.DeleteProcessesIdResponse {

	processRep := repository.NewProcessRepository()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		panic(err)
	}

	processRep.Delete(int32(id))

	return openapi_models.DeleteProcessesIdResponse{}

}
