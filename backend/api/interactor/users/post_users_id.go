package users

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/kenkonno/gantt-chart-proto/backend/api/constants"
	"github.com/kenkonno/gantt-chart-proto/backend/api/middleware"
	"github.com/kenkonno/gantt-chart-proto/backend/api/openapi_models"
	"github.com/kenkonno/gantt-chart-proto/backend/api/utils"
	"github.com/kenkonno/gantt-chart-proto/backend/models/db"
	"github.com/kenkonno/gantt-chart-proto/backend/repository"
	"net/http"
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

	// Role.Workerの時のみ自分の更新かのチェックをする
	sessionUser := userRep.Find(*middleware.GetUserId(c))
	if sessionUser.Role == constants.RoleWorker && *userReq.User.Id != *sessionUser.Id {
		c.JSON(http.StatusForbidden, "他者を更新する権限がありません。")
		return openapi_models.PostUsersIdResponse{}, errors.New("forbidden")
	}

	// パスワードをハッシュ化
	pw, err := utils.CryptPassword(userReq.User.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return openapi_models.PostUsersIdResponse{}, errors.New("パスワードの暗号化に失敗しました。")
	}
	oldUser := userRep.Find(*userReq.User.Id)

	// 登録済みのEmailははじく
	userIsExists := userRep.FindByEmail(userReq.User.Email)
	// メールアドレスの重複チェック。メールアドレスが変更されているかつ既に存在していればNG
	if oldUser.Email != userReq.User.Email && userIsExists.Id != nil {
		c.JSON(http.StatusBadRequest, "メールアドレスが重複しています。")
		return openapi_models.PostUsersIdResponse{}, errors.New("Email is already taken.")
	}

	// パスワードリセット 空文字の時は前回の設定を引き継ぐ（管理者が別のユーザーを更新するケースがあるため）
	passwordReset := false
	// パスワードリセット済みかつ、パスワードは空文字の場合は更新しない。
	if userReq.User.Password == "" && oldUser.PasswordReset == true {
		pw = oldUser.Password
		passwordReset = oldUser.PasswordReset
	} else {
		// バリデーション
		if !utils.ValidatePassword(userReq.User.Password) {
			// TODO: 本来はカスタムエラーを作ってエラーハンドリングすべきだが速度優先で一旦エラーを返したらapi側ではレスポンスの制御をしないようにした。
			c.JSON(http.StatusBadRequest, constants.E0001)
			return openapi_models.PostUsersIdResponse{}, errors.New("")
		}
		passwordReset = true
	}

	userRep.Upsert(db.User{
		Id:                  userReq.User.Id,
		DepartmentId:        userReq.User.DepartmentId,
		LimitOfOperation:    userReq.User.LimitOfOperation,
		LastName:            strings.TrimSpace(userReq.User.LastName),
		FirstName:           strings.TrimSpace(userReq.User.FirstName),
		Password:            pw,
		Email:               strings.ToLower(strings.TrimSpace(userReq.User.Email)),
		Role:                userReq.User.Role,
		PasswordReset:       passwordReset,
		EmploymentStartDate: userReq.User.EmploymentStartDate,
		EmploymentEndDate:   userReq.User.EmploymentEndDate,
		CreatedAt:           time.Time{},
		UpdatedAt:           0,
	})

	return openapi_models.PostUsersIdResponse{}, nil

}
