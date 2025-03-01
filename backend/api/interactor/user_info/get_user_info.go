package user_info

import (
	"github.com/gin-gonic/gin"
	"github.com/kenkonno/gantt-chart-proto/backend/api/constants"
	"github.com/kenkonno/gantt-chart-proto/backend/api/middleware"
	"github.com/kenkonno/gantt-chart-proto/backend/api/openapi_models"
	"github.com/kenkonno/gantt-chart-proto/backend/repository"
	"github.com/kenkonno/gantt-chart-proto/backend/repository/simulation"
)

// Tokenからユーザー情報を返却する
func GetUserInfoInvoke(c *gin.Context) (openapi_models.GetUserInfoResponse, error) {
	userId := middleware.GetUserId(c)
	// セッション切れの場合は空で戻す
	if userId == nil {
		return openapi_models.GetUserInfoResponse{}, nil
	}

	userRep := repository.NewUserRepository(middleware.GetRepositoryMode(c)...)
	simulationLockRep := simulation.NewSimulationLock()
	simulationLock := simulationLockRep.Find(constants.SimulateTypeSchedule)

	isSimulateUser := simulationLock.Status == constants.SimulateStatusInProgress && simulationLock.LockedBy == *userId

	var userInfoResponse openapi_models.GetUserInfoResponse
	user := userRep.Find(*userId)
	userInfoResponse = openapi_models.GetUserInfoResponse{
		User: openapi_models.User{
			Id:                  user.Id,
			DepartmentId:        user.DepartmentId,
			LimitOfOperation:    user.LimitOfOperation,
			LastName:            user.LastName,
			FirstName:           user.FirstName,
			Password:            "",
			Email:               user.Email,
			CreatedAt:           user.CreatedAt,
			UpdatedAt:           user.UpdatedAt,
			Role:                user.Role,
			PasswordReset:       user.PasswordReset,
			EmploymentStartDate: user.EmploymentStartDate,
			EmploymentEndDate:   user.EmploymentEndDate,
		},
		IsSimulateUser: isSimulateUser,
	}
	return userInfoResponse, nil
}
