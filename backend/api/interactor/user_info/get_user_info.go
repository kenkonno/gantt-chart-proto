package user_info

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/kenkonno/gantt-chart-proto/backend/api/middleware"
	"github.com/kenkonno/gantt-chart-proto/backend/api/openapi_models"
	"github.com/kenkonno/gantt-chart-proto/backend/repository"
	"strconv"
)

// Tokenからユーザー情報を返却する
func GetUserInfoInvoke(c *gin.Context) openapi_models.GetUserInfoResponse {
	userRep := repository.NewUserRepository()
	sessionID, err := c.Cookie("session_id")
	fmt.Println(1)
	if err != nil {
		return openapi_models.GetUserInfoResponse{}
	}
	fmt.Println(2)
	strUserId := middleware.GetUserId(sessionID)
	if strUserId == nil {
		return openapi_models.GetUserInfoResponse{}
	}
	fmt.Println(3)
	int32UserId, _ := strconv.ParseInt(*strUserId, 10, 32)
	userId := int32(int32UserId)
	user := userRep.Find(userId)
	fmt.Println(4)
	return openapi_models.GetUserInfoResponse{
		User: openapi_models.User{
			Id:               user.Id,
			DepartmentId:     user.DepartmentId,
			LimitOfOperation: user.LimitOfOperation,
			Name:             user.Name,
			Password:         "",
			Email:            user.Email,
			CreatedAt:        user.CreatedAt,
			UpdatedAt:        user.UpdatedAt,
			Role:             user.Role,
		},
	}
}