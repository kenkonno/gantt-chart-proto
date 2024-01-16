package users

import (
	"github.com/gin-gonic/gin"
	"github.com/kenkonno/gantt-chart-proto/backend/api/openapi_models"
	"github.com/kenkonno/gantt-chart-proto/backend/models/db"
	"github.com/kenkonno/gantt-chart-proto/backend/repository"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

func PostUsersIdInvoke(c *gin.Context) openapi_models.PostUsersIdResponse {

	userRep := repository.NewUserRepository()

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
	// パスワードは空文字の場合は更新しない。
	if userReq.User.Password == "" {
		pw = oldUser.Password
	}

	userRep.Upsert(db.User{
		Id:               userReq.User.Id,
		DepartmentId:     userReq.User.DepartmentId,
		LimitOfOperation: userReq.User.LimitOfOperation,
		Name:             userReq.User.Name,
		Password:         pw,
		Email:            userReq.User.Email,
		Role:             userReq.User.Role,
		CreatedAt:        time.Time{},
		UpdatedAt:        0,
	})

	return openapi_models.PostUsersIdResponse{}

}
