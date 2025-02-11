package login

import (
	"fmt"
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

func PostLoginInvoke(c *gin.Context) (openapi_models.PostLoginResponse, error) {

	var userReq openapi_models.PostLoginRequest
	if err := c.ShouldBindJSON(&userReq); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		panic(err)
	}

	/**
	NOTE: かなり微妙な設計。ゲストから通常ログインの復帰のためのコード。
	権限チェックのmiddlewareで既にゲストモードとして特定されているが、repository_modeを上書きして通常モードにしている。
	ただし、シミュレーション実行中ユーザーのことを考えるとおかしくなりそう。一旦通常テーブルを見に行けば運用上問題なさそうなので、こうする。
	 */
	if userReq.Id != "" && userReq.Password != "" {
		c.Set("repository_mode", []string{})
	}

	userRep := repository.NewUserRepository(middleware.GetRepositoryMode(c)...)
	facilitySharedLinkRep := repository.NewFacilitySharedLinkRepository(middleware.GetRepositoryMode(c)...)


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
	fmt.Println("############################ ", isAuthenticated, middleware.GetRepositoryMode(c))

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

	return openapi_models.PostLoginResponse{}, nil

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
