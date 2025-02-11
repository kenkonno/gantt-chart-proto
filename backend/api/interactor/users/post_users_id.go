package users

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/kenkonno/gantt-chart-proto/backend/api/middleware"
	"github.com/kenkonno/gantt-chart-proto/backend/api/openapi_models"
	"github.com/kenkonno/gantt-chart-proto/backend/models/db"
	"github.com/kenkonno/gantt-chart-proto/backend/repository"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"regexp"
	"strings"
	"time"
)

func PostUsersIdInvoke(c *gin.Context) (openapi_models.PostUsersIdResponse, error) {

	userRep := repository.NewUserRepository(middleware.GetRepositoryMode(c)...)

	var userReq openapi_models.PostUsersRequest
	if err := c.ShouldBindJSON(&userReq); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		panic(err)
	}
	// パスワードをハッシュ化
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userReq.User.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		panic(err)
	}
	oldUser := userRep.Find(*userReq.User.Id)
	pw := string(hashedPassword)

	// パスワードリセット
	passwordReset := false
	if oldUser.PasswordReset {
		// 過去設定済みであればtrue固定
		passwordReset = true
	} else {
		// パスワードが更新されていない場合はエラーとする
		if !validatePassword(userReq.User.Password) {
			// TODO: 本来はカスタムエラーを作ってエラーハンドリングすべきだが速度優先で一旦エラーを返したらapi側ではレスポンスの制御をしないようにした。
			c.JSON(http.StatusBadRequest, `パスワードは小文字、大文字、数字、特殊文字「 [!@#\$%^&*()] 」を含み8文字以上にしてください。`)
			return openapi_models.PostUsersIdResponse{}, errors.New("")
		}
		passwordReset = true
	}

	// パスワードは空文字の場合は更新しない。
	if userReq.User.Password == "" {
		pw = oldUser.Password
	}

	userRep.Upsert(db.User{
		Id:               userReq.User.Id,
		DepartmentId:     userReq.User.DepartmentId,
		LimitOfOperation: userReq.User.LimitOfOperation,
		LastName:         strings.TrimSpace(userReq.User.LastName),
		FirstName:        strings.TrimSpace(userReq.User.FirstName),
		Password:         pw,
		Email:            strings.TrimSpace(userReq.User.Email),
		Role:             userReq.User.Role,
		CreatedAt:        time.Time{},
		UpdatedAt:        0,
		PasswordReset:    passwordReset,
	})

	return openapi_models.PostUsersIdResponse{}, nil

}

func validatePassword(password string) bool {
	var (
		containsMin   = regexp.MustCompile(`[a-z]`).MatchString
		containsMax   = regexp.MustCompile(`[A-Z]`).MatchString
		containsNum   = regexp.MustCompile(`[0-9]`).MatchString
		containsSpecial = regexp.MustCompile(`[!@#\$%^&*()]`).MatchString
		lengthValid   = regexp.MustCompile(`.{8,}`).MatchString
	)

	return containsMin(password) && containsMax(password) && containsNum(password) && containsSpecial(password) && lengthValid(password)
}