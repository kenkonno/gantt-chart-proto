package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/kenkonno/gantt-chart-proto/backend/api/constants"
	"github.com/kenkonno/gantt-chart-proto/backend/repository"
	"net/http"
	"strconv"
)

func GuestAccessControl() gin.HandlerFunc {
	return func(c *gin.Context) {

		// ゲストか確認
		if !IsGuest(c) {
			c.Next()
			return
		}

		// ゲストの場合は GetのリクエストパラメータでfacilityIdはuuidに紐づくものだけ利用可能にする
		uuid, err := c.Cookie(constants.FacilitySharedLinkUUID)
		if err != nil {
			// ゲストでUUIDが無ければエラーとする。
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "no permission to this resource"})
			return
		}

		// facilityIdが関連しないリクエストならほかの権限に任せる。
		facilityIdStr := c.Query("facilityId")
		if facilityIdStr == "" {
			c.Next()
			return
		}
		facilityId, _ := strconv.Atoi(facilityIdStr)

		// facilityIdが関連するリクエストなら整合性チェックをする。
		facilitySharedLinkRep := repository.NewFacilitySharedLinkRepository()
		facilitySharedLink := facilitySharedLinkRep.FindByUUID(uuid)

		if facilitySharedLink.FacilityId != int32(facilityId) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "no permission to this resource"})
			return
		}
		c.Next()
	}
}
