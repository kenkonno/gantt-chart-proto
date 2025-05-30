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
	"strings"
	"time"
)

func PostUsersInvoke(c *gin.Context) (openapi_models.PostUsersResponse, error) {

	userRep := repository.NewUserRepository(middleware.GetRepositoryMode(c)...)

	var userReq openapi_models.PostUsersRequest
	if err := c.ShouldBindJSON(&userReq); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		panic(err)
	}

	// 登録済みのEmailははじく
	userIsExists := userRep.FindByEmail(userReq.User.Email)
	if userIsExists.Id != nil {
		c.JSON(http.StatusBadRequest, "登録済みのメールアドレスです。")
		return openapi_models.PostUsersResponse{}, errors.New("Email is already taken.")
	}

	// パスワードをハッシュ化
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userReq.User.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		panic(err)
	}
	userRep.Upsert(db.User{
		DepartmentId:        userReq.User.DepartmentId,
		LimitOfOperation:    userReq.User.LimitOfOperation,
		LastName:            strings.TrimSpace(userReq.User.LastName),
		FirstName:           strings.TrimSpace(userReq.User.FirstName),
		Password:            string(hashedPassword),
		Email:               strings.ToLower(strings.TrimSpace(userReq.User.Email)),
		Role:                userReq.User.Role,
		CreatedAt:           time.Time{},
		UpdatedAt:           0,
		EmploymentStartDate: userReq.User.EmploymentStartDate,
		EmploymentEndDate:   userReq.User.EmploymentEndDate,
	})

	return openapi_models.PostUsersResponse{}, nil

}
