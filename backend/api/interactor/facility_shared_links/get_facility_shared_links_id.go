package facility_shared_links

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/kenkonno/gantt-chart-proto/backend/api/middleware"
	"github.com/kenkonno/gantt-chart-proto/backend/api/openapi_models"
	"github.com/kenkonno/gantt-chart-proto/backend/models/db"
	"github.com/kenkonno/gantt-chart-proto/backend/repository"
)

// facilityIdをもとに取得する
func GetFacilitySharedLinksIdInvoke(c *gin.Context) (openapi_models.GetFacilitySharedLinksIdResponse, error) {
	facilitySharedLinkRep := repository.NewFacilitySharedLinkRepository(middleware.GetRepositoryMode(c)...)

	// REST的にはおかしいが、工数簡略化のためidパラメータをfacilityIdとして利用する。
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		panic(err)
	}

	facilitySharedLink := facilitySharedLinkRep.FindByFacilityId(int32(id))

	// 存在しない場合は新規作成する
	if facilitySharedLink.Uuid == "" {
		facilitySharedLinkRep.Upsert(db.FacilitySharedLink{
			Id:         nil,
			FacilityId: int32(id),
			Uuid:       uuid.New().String(),
			CreatedAt:  time.Time{},
			UpdatedAt:  0,
		})
		facilitySharedLink = facilitySharedLinkRep.FindByFacilityId(int32(id))
	}

	return openapi_models.GetFacilitySharedLinksIdResponse{
		FacilitySharedLink: openapi_models.FacilitySharedLink{
			Id:         facilitySharedLink.Id,
			FacilityId: facilitySharedLink.FacilityId,
			Uuid:       &facilitySharedLink.Uuid,
			CreatedAt:  facilitySharedLink.CreatedAt,
			UpdatedAt: int32(facilitySharedLink.UpdatedAt),
		},
	}, nil
}
