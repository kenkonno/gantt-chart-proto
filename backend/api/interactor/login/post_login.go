package login

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/kenkonno/gantt-chart-proto/backend/api/constants"
	"github.com/kenkonno/gantt-chart-proto/backend/api/middleware"
	"github.com/kenkonno/gantt-chart-proto/backend/api/openapi_models"
	"github.com/kenkonno/gantt-chart-proto/backend/repository"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"os"
)

func PostLoginInvoke(c *gin.Context) openapi_models.PostLoginResponse {

	userRep := repository.NewUserRepository()
	facilitySharedLinkRep := repository.NewFacilitySharedLinkRepository()

	var userReq openapi_models.PostLoginRequest
	if err := c.ShouldBindJSON(&userReq); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		panic(err)
	}

	isAuthenticated := false
	var userId *int32

	// ゲストログインの処理を実行する
	if userReq.Uuid != "" {
		facilitySharedLink := facilitySharedLinkRep.FindByUUID(userReq.Uuid)
		if facilitySharedLink.Id != nil {
			isAuthenticated = true
			guestId := int32(constants.GuestID)
			userId = &guestId
			writeCookie(constants.FacilitySharedLinkUUID, userReq.Uuid, c)
		}
	} else {
		// 通常ログインの処理を実行する
		user := userRep.FindByEmail(userReq.Id) // IDといいつつEmail
		isAuthenticated = VerifyPassword(userReq.Password, user.Password)
		userId = user.Id
	}

	if isAuthenticated {
		sessionId := uuid.New().String()
		// セッションを作成する
		name := "session_id"
		value := sessionId
		// Set the cookie
		writeCookie(name, value, c)
		// redisに書き込む
		middleware.UpdateSessionID(sessionId, *userId)
	} else {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
	}

	return openapi_models.PostLoginResponse{}

}

func VerifyPassword(inputPassword string, hashedPassword string) bool {
	// ハッシュと入力されたパスワードを比較
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(inputPassword))

	// エラーがない場合、パスワードは一致している
	return err == nil
}

func writeCookie(name string, value string, c *gin.Context) {
	maxAge := 86400
	path := "/"
	domain := os.Getenv("HOST_NAME")
	secure := false
	httpOnly := true
	// Set the cookie
	c.SetCookie(name, value, maxAge, path, domain, secure, httpOnly)
}
