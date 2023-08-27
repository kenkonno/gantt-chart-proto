package users

import (
	"github.com/gin-gonic/gin"
	"github.com/kenkonno/gantt-chart-proto/backend/api/openapi_models"
	"github.com/kenkonno/gantt-chart-proto/backend/models/db"
	"github.com/kenkonno/gantt-chart-proto/backend/repository"
	"github.com/samber/lo"
)

func GetUsersInvoke(c *gin.Context) openapi_models.GetUsersResponse {
	userRep := repository.NewUserRepository()

	userList := userRep.FindAll()

	return openapi_models.GetUsersResponse{
		List: lo.Map(userList, func(item db.User, index int) openapi_models.User {
			return openapi_models.User{
				Id:        item.Id,
				Name:      item.Name,
				Password:  item.Password,
				Email:     item.Email,
				CreatedAt: item.CreatedAt,
				UpdatedAt: item.UpdatedAt,
			}
		}),
	}
}
