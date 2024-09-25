package simulation

import (
	"github.com/gin-gonic/gin"
	"github.com/kenkonno/gantt-chart-proto/backend/api/constants"
	"github.com/kenkonno/gantt-chart-proto/backend/api/middleware"
	"github.com/kenkonno/gantt-chart-proto/backend/api/openapi_models"
	"github.com/kenkonno/gantt-chart-proto/backend/models/db"
	"github.com/kenkonno/gantt-chart-proto/backend/repository"
	"github.com/kenkonno/gantt-chart-proto/backend/repository/common"
	"github.com/samber/lo"
)

func GetSimulationInvoke(c *gin.Context) openapi_models.GetSimulationResponse {
	facilityRep := repository.NewFacilityRepository()

	facilityList := facilityRep.FindAll([]string{}, []string{})

	// ゲストユーザーの場合はUUIDに紐づく設備のみ返す
	// TODO: 本来ならば FindAll をそのようにすべき。Mock的な感じでDIするようにする。
	if middleware.IsGuest(c) {
		facilitySharedLinkRep := common.NewFacilitySharedLinkRepository()
		uuid, _ := c.Cookie(constants.FacilitySharedLinkUUID)
		facilitySharedLink := facilitySharedLinkRep.FindByUUID(uuid)
		facilityList = lo.Filter(facilityList, func( item db.Facility, index int) bool {
			return facilitySharedLink.FacilityId == *item.Id
		})
	}

	return openapi_models.GetSimulationResponse{
	}
}
