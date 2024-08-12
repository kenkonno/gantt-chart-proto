package logout

import (
	"github.com/gin-gonic/gin"
	"github.com/kenkonno/gantt-chart-proto/backend/api/middleware"
	"github.com/kenkonno/gantt-chart-proto/backend/api/openapi_models"
)

func PostLogoutInvoke(c *gin.Context) openapi_models.PostLogoutResponse {

	// セッションとクッキーをクリアする
	sessionId, _ := c.Cookie("session_id")
	if sessionId == "" {
		return openapi_models.PostLogoutResponse{}
	}
	middleware.ClearSession(sessionId)
	return openapi_models.PostLogoutResponse{}

}
