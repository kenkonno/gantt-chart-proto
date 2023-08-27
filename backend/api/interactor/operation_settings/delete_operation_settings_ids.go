package operation_settings

import (
	"github.com/gin-gonic/gin"
	"github.com/kenkonno/gantt-chart-proto/backend/api/openapi_models"
	"github.com/kenkonno/gantt-chart-proto/backend/repository"
	"strconv"
)

func DeleteOperationSettingsIdInvoke(c *gin.Context) openapi_models.DeleteOperationSettingsIdResponse {

	operationSettingRep := repository.NewOperationSettingRepository()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		panic(err)
	}

	operationSettingRep.Delete(int32(id))

	return openapi_models.DeleteOperationSettingsIdResponse{}

}
