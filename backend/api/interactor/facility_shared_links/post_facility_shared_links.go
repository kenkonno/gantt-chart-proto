package facility_shared_links

import (
	"github.com/gin-gonic/gin"
	"github.com/kenkonno/gantt-chart-proto/backend/api/middleware"
	"github.com/kenkonno/gantt-chart-proto/backend/api/openapi_models"
	"github.com/kenkonno/gantt-chart-proto/backend/models/db"
	"github.com/kenkonno/gantt-chart-proto/backend/repository"
	"net/http"
	"time"
)

// 未使用。取得時に存在しなければ追加するため。
func PostFacilitySharedLinksInvoke(c *gin.Context) (openapi_models.PostFacilitySharedLinksResponse, error) {

	facilitySharedLinkRep := repository.NewFacilitySharedLinkRepository(middleware.GetRepositoryMode(c)...)

	var facilitySharedLinkReq openapi_models.PostFacilitySharedLinksRequest
	if err := c.ShouldBindJSON(&facilitySharedLinkReq); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		panic(err)
	}
	facilitySharedLinkRep.Upsert(db.FacilitySharedLink{
		FacilityId: facilitySharedLinkReq.FacilitySharedLink.FacilityId,
		Uuid:       *facilitySharedLinkReq.FacilitySharedLink.Uuid,
		CreatedAt:  time.Time{},
		UpdatedAt:  0,
	})

	return openapi_models.PostFacilitySharedLinksResponse{}, nil

}
