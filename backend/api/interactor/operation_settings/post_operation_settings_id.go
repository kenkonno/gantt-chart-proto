package operation_settings

import (
	"github.com/gin-gonic/gin"
	"github.com/kenkonno/gantt-chart-proto/backend/api/middleware"
	"github.com/kenkonno/gantt-chart-proto/backend/api/openapi_models"
	"github.com/kenkonno/gantt-chart-proto/backend/models/db"
	"github.com/kenkonno/gantt-chart-proto/backend/repository"
	"net/http"
	"strconv"
)

func PostOperationSettingsIdInvoke(c *gin.Context) (openapi_models.PostOperationSettingsIdResponse, error) {

	operationSettingRep := repository.NewOperationSettingRepository(middleware.GetRepositoryMode(c)...)
	facilityId, err := strconv.Atoi(c.Param("id")) // facility_id
	if err != nil {
		panic(err)
	}

	var operationSettingReq openapi_models.PostOperationSettingsRequest
	if err := c.ShouldBindJSON(&operationSettingReq); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		panic(err)
	}

	for _, v := range operationSettingReq.OperationSettings {
		for _, vv := range v.WorkHours {
			operationSettingRep.Upsert(db.OperationSetting{
				Id:         v.Id,
				FacilityId: int32(facilityId),
				UnitId:     v.UnitId,
				ProcessId:  vv.ProcessId,
				WorkHour:   vv.WorkHour,
			})
		}
	}

	return openapi_models.PostOperationSettingsIdResponse{}, nil
	//
}
