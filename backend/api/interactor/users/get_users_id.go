package users

import (
	"github.com/gin-gonic/gin"
	"github.com/kenkonno/gantt-chart-proto/backend/api/middleware"
	"github.com/kenkonno/gantt-chart-proto/backend/api/openapi_models"
	"github.com/kenkonno/gantt-chart-proto/backend/repository"
	"strconv"
)

func GetUsersIdInvoke(c *gin.Context) (openapi_models.GetUsersIdResponse, error) {
	userRep := repository.NewUserRepository(middleware.GetRepositoryMode(c)...)

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		panic(err)
	}

	user := userRep.Find(int32(id))

	return openapi_models.GetUsersIdResponse{
		User: openapi_models.User{
			Id:                  user.Id,
			DepartmentId:        user.DepartmentId,
			LimitOfOperation:    user.LimitOfOperation,
			LastName:            user.LastName,
			FirstName:           user.FirstName,
			Password:            "", // Passwordはユーザーに含めない
			Email:               user.Email,
			Role:                user.Role,
			CreatedAt:           user.CreatedAt,
			UpdatedAt:           user.UpdatedAt,
			EmploymentStartDate: user.EmploymentStartDate,
			EmploymentEndDate:   user.EmploymentEndDate,
		},
	}, nil
}
