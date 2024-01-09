package login

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/kenkonno/gantt-chart-proto/backend/api/middleware"
	"github.com/kenkonno/gantt-chart-proto/backend/api/openapi_models"
	"github.com/kenkonno/gantt-chart-proto/backend/repository"
	"net/http"
	"os"
)

func PostLoginInvoke(c *gin.Context) openapi_models.PostLoginResponse {

	userRep := repository.NewUserRepository()

	var userReq openapi_models.PostLoginRequest
	if err := c.ShouldBindJSON(&userReq); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		panic(err)
	}

	user := userRep.FindByAuth(userReq.Id, userReq.Password)
	if user.Id == nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
	} else {
		sessionId := uuid.New().String()
		// セッションを作成する
		name := "session_id"
		value := sessionId
		maxAge := 86400
		path := "/"
		domain := os.Getenv("HOST_NAME")
		secure := false
		httpOnly := true
		// Set the cookie
		c.SetCookie(name, value, maxAge, path, domain, secure, httpOnly)
		// redisに書き込む
		middleware.UpdateSessionID(sessionId, *user.Id)
	}

	return openapi_models.PostLoginResponse{}

}
