package login

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/kenkonno/gantt-chart-proto/backend/api/middleware"
	"github.com/kenkonno/gantt-chart-proto/backend/api/openapi_models"
	"github.com/kenkonno/gantt-chart-proto/backend/repository"
	"golang.org/x/crypto/bcrypt"
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
	user := userRep.FindByEmail(userReq.Id) // IDといいつつEmail
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userReq.Password), bcrypt.DefaultCost)
	fmt.Println(string(hashedPassword), err)

	if !VerifyPassword(userReq.Password, user.Password) {
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


func VerifyPassword(inputPassword string, hashedPassword string) bool {
	// ハッシュと入力されたパスワードを比較
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(inputPassword))

	// エラーがない場合、パスワードは一致している
	return err == nil
}