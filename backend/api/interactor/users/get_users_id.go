package users

import (
	"github.com/gin-gonic/gin"
	"github.com/kenkonno/gantt-chart-proto/backend/api/constants"
	"github.com/kenkonno/gantt-chart-proto/backend/api/middleware"
	"github.com/kenkonno/gantt-chart-proto/backend/api/openapi_models"
	"github.com/kenkonno/gantt-chart-proto/backend/repository"
	"strconv"
)

func GetUsersIdInvoke(c *gin.Context) openapi_models.GetUsersIdResponse {
	sessionID, err := c.Cookie("session_id")
	// TODO: この辺のDIをAPIとinteractorで分けて処理するべきだが、自動生成の兼ね合いで対応できず。コスト的にはinteractorに記述したほうがいったんはよい？と思ったけどそうでもないか。GetInteractorで gin.Contextを渡してInteractorの構造を返す形にするのがいいかも。そうするとロジックとDIで分離できる。シミュレーション機能のリファクタリングのタイミングで実施する。
	userRep := repository.NewUserRepository()
	if *middleware.GetUserId(sessionID) == constants.GuestID {
		userRep = repository.NewUserRepository(repository.GuestMode)
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		panic(err)
	}

	user := userRep.Find(int32(id))

	return openapi_models.GetUsersIdResponse{
		User: openapi_models.User{
			Id:               user.Id,
			DepartmentId:     user.DepartmentId,
			LimitOfOperation: user.LimitOfOperation,
			LastName:         user.LastName,
			FirstName:        user.FirstName,
			Password:         "", // Passwordはユーザーに含めない
			Email:            user.Email,
			Role:             user.Role,
			CreatedAt:        user.CreatedAt,
			UpdatedAt:        user.UpdatedAt,
		},
	}
}
