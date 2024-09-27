package facility_shared_links

import (
	"github.com/gin-gonic/gin"
	"github.com/kenkonno/gantt-chart-proto/backend/api/middleware"
	"github.com/kenkonno/gantt-chart-proto/backend/api/openapi_models"
	"github.com/kenkonno/gantt-chart-proto/backend/models/db"
	"github.com/kenkonno/gantt-chart-proto/backend/repository"
	"github.com/samber/lo"
)

// 未使用
func GetFacilitySharedLinksInvoke(c *gin.Context) openapi_models.GetFacilitySharedLinksResponse {
	facilitySharedLinkRep := repository.NewFacilitySharedLinkRepository(middleware.GetRepositoryMode(c)...)

	facilitySharedLinkList := facilitySharedLinkRep.FindAll()

	return openapi_models.GetFacilitySharedLinksResponse{
		List: lo.Map(facilitySharedLinkList, func(item db.FacilitySharedLink, index int) openapi_models.FacilitySharedLink {
			return openapi_models.FacilitySharedLink{
				Id:         item.Id,
				FacilityId: item.FacilityId,
				Uuid:       &item.Uuid,
				CreatedAt:  item.CreatedAt,
				UpdatedAt:  int(item.UpdatedAt),
			}
		}),
	}
}
