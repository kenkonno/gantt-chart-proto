package facilities

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

func GetFacilitiesInvoke(c *gin.Context) openapi_models.GetFacilitiesResponse {
	facilityRep := repository.NewFacilityRepository(middleware.GetRepositoryMode(c)...)

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

	return openapi_models.GetFacilitiesResponse{
		List: lo.Map(facilityList, func(item db.Facility, index int) openapi_models.Facility {
			return openapi_models.Facility{
				Id:              item.Id,
				Name:            item.Name,
				TermFrom:        item.TermFrom,
				TermTo:          item.TermTo,
				Order:           int32(item.Order),
				CreatedAt:       item.CreatedAt,
				UpdatedAt:       item.UpdatedAt,
				Status:          item.Status,
				Type:            item.Type,
				ShipmentDueDate: item.ShipmentDueDate,
			}
		}),
	}
}
