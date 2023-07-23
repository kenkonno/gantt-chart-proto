package users

import (
	"github.com/gin-gonic/gin"
	"github.com/kenkonno/gantt-chart-proto/backend/api/openapi_models"
	"github.com/kenkonno/gantt-chart-proto/backend/models/db"
	"github.com/kenkonno/gantt-chart-proto/backend/repository"
)

func PostUsersInvoke(c *gin.Context) openapi_models.PostUsersResponse {

	userRep := repository.NewUserRepository()

	var userReq openapi_models.PostUsersRequest
	if err := c.ShouldBindJSON(&userReq); err != nil {
		panic("invalid json")
	}
	userRep.Upsert(db.User{
		Id:        userReq.User.Id,
		Password:  userReq.User.Password,
		Email:     userReq.User.Email,
		CreatedAt: userReq.User.CreatedAt,
		UpdatedAt: userReq.User.UpdatedAt,
	})

	return openapi_models.PostUsersResponse{}

}
