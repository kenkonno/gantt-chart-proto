package operation_settings
import (
	"github.com/gin-gonic/gin"
	"github.com/kenkonno/gantt-chart-proto/backend/api/openapi_models"
	"github.com/kenkonno/gantt-chart-proto/backend/repository"
)
func PostOperationSettingsIdInvoke(c *gin.Context) openapi_models.PostOperationSettingsIdResponse {

	operationSettingRep := repository.NewOperationSettingRepository()

	var operationSettingReq openapi_models.PostOperationSettingsRequest
	if err := c.ShouldBindJSON(&operationSettingReq); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		panic(err)
	}

	operationSettingRep.Upsert(db.OperationSetting{
				Id:        operationSettingReq.OperationSetting.Id,
				FacilityId:        operationSettingReq.OperationSetting.FacilityId,
				UnitId:        operationSettingReq.OperationSetting.UnitId,
				ProcessId:        operationSettingReq.OperationSetting.ProcessId,
				WorkHour:        operationSettingReq.OperationSetting.WorkHour,
				CreatedAt:        time.Time{},
				UpdatedAt:        0,

	})

	return openapi_models.PostOperationSettingsIdResponse{}

}

