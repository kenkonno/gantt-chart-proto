package users

import (
	"github.com/gin-gonic/gin"
	"github.com/kenkonno/gantt-chart-proto/backend/api/middleware"
	"github.com/kenkonno/gantt-chart-proto/backend/api/openapi_models"
	"github.com/kenkonno/gantt-chart-proto/backend/models/db"
	"github.com/kenkonno/gantt-chart-proto/backend/repository"
	"github.com/samber/lo"
)

func GetUsersInvoke(c *gin.Context) (openapi_models.GetUsersResponse, error) {
	userRep := repository.NewUserRepository(middleware.GetRepositoryMode(c)...)

	userList := userRep.FindAll()

	return openapi_models.GetUsersResponse{
		List: lo.Map(userList, func(item db.User, index int) openapi_models.User {
			return openapi_models.User{
				Id:                  item.Id,
				DepartmentId:        item.DepartmentId,
				LimitOfOperation:    item.LimitOfOperation,
				LastName:            item.LastName,
				FirstName:           item.FirstName,
				Password:            "", // Passwordはレスポンスに含めない
				Email:               item.Email,
				CreatedAt:           item.CreatedAt,
				UpdatedAt:           item.UpdatedAt,
				Role:                item.Role,
				PasswordReset:       item.PasswordReset,
				EmploymentStartDate: item.EmploymentStartDate,
				EmploymentEndDate:   item.EmploymentEndDate,
			}
		}),
	}, nil
}
