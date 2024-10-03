package operation_settings

import (
	"github.com/gin-gonic/gin"
	"github.com/kenkonno/gantt-chart-proto/backend/api/middleware"
	"github.com/kenkonno/gantt-chart-proto/backend/api/openapi_models"
	"github.com/kenkonno/gantt-chart-proto/backend/repository"
	"strconv"
)

func GetOperationSettingsIdInvoke(c *gin.Context) openapi_models.GetOperationSettingsIdResponse {
	operationSettingRep := repository.NewOperationSettingRepository(middleware.GetRepositoryMode(c)...)

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		panic(err)
	}

	operationSettings := operationSettingRep.FindByFacilityId(int32(id))

	var results = []openapi_models.OperationSetting{}

	if len(operationSettings) > 0 {
		var prev = operationSettings[0]
		var workHours []openapi_models.WorkHour
		// 工程の集約
		for _, v := range operationSettings {
			if prev.UnitId != v.UnitId {
				results = append(results, openapi_models.OperationSetting{
					Id:         prev.Id,
					FacilityId: prev.FacilityId,
					UnitId:     prev.UnitId,
					WorkHours:  append([]openapi_models.WorkHour{}, workHours...), // copy slice
					CreatedAt:  prev.CreatedAt,
					UpdatedAt:  prev.UpdatedAt,
				})
				workHours = []openapi_models.WorkHour{}
			}
			workHours = append(workHours, openapi_models.WorkHour{
				ProcessId: v.ProcessId,
				WorkHour:  v.WorkHour,
			})
			prev = v
		}
		// 最終行の処理
		results = append(results, openapi_models.OperationSetting{
			Id:         prev.Id,
			FacilityId: prev.FacilityId,
			UnitId:     prev.UnitId,
			WorkHours:  append([]openapi_models.WorkHour{}, workHours...), // copy slice
			CreatedAt:  prev.CreatedAt,
			UpdatedAt:  0,
		})
	}

	return openapi_models.GetOperationSettingsIdResponse{
		OperationSettings: results,
	}
}
