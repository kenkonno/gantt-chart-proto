package users

import (
	"github.com/gin-gonic/gin"
	"github.com/kenkonno/gantt-chart-proto/backend/api/openapi_models"
	"github.com/kenkonno/gantt-chart-proto/backend/repository"
	"strconv"
)

func GetUsersIdInvoke(c *gin.Context) openapi_models.GetUsersIdResponse {
	userRep := repository.NewUserRepository()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		panic(err)
	}

	user := userRep.Find(int32(id))

	return openapi_models.GetUsersIdResponse{
		User: openapi_models.User{
			Id:        user.Id,
			Password:  user.Password,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		},
	}
}
