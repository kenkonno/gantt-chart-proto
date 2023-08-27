package operation_settings
import (
	"github.com/gin-gonic/gin"
	"github.com/kenkonno/gantt-chart-proto/backend/api/openapi_models"
	"github.com/kenkonno/gantt-chart-proto/backend/repository"
)
func GetOperationSettingsIdInvoke(c *gin.Context) openapi_models.GetOperationSettingsIdResponse {
	operationSettingRep := repository.NewOperationSettingRepository()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		panic(err)
	}

	operationSetting := operationSettingRep.Find(int32(id))

	return openapi_models.GetOperationSettingResponse{
		OperationSetting: openapi_models.OperationSetting{
				Id:        operationSetting.Id,
				FacilityId:        operationSetting.FacilityId,
				UnitId:        operationSetting.UnitId,
				ProcessId:        operationSetting.ProcessId,
				WorkHour:        operationSetting.WorkHour,
				CreatedAt:        operationSetting.CreatedAt,
				UpdatedAt:        operationSetting.UpdatedAt,

		},
	}
}
