package user_info

import (
	"github.com/gin-gonic/gin"
	"github.com/kenkonno/gantt-chart-proto/backend/api/constants"
	"github.com/kenkonno/gantt-chart-proto/backend/api/middleware"
	"github.com/kenkonno/gantt-chart-proto/backend/api/openapi_models"
	"github.com/kenkonno/gantt-chart-proto/backend/repository"
	"time"
)

// Tokenからユーザー情報を返却する
func GetUserInfoInvoke(c *gin.Context) openapi_models.GetUserInfoResponse {
	sessionID, err := c.Cookie("session_id")
	if err != nil {
		return openapi_models.GetUserInfoResponse{}
	}
	// TODO: この辺のDIをAPIとinteractorで分けて処理するべきだが、自動生成の兼ね合いで対応できず。コスト的にはinteractorに記述したほうがいったんはよい？と思ったけどそうでもないか。GetInteractorで gin.Contextを渡してInteractorの構造を返す形にするのがいいかも。そうするとロジックとDIで分離できる。シミュレーション機能のリファクタリングのタイミングで実施する。
	userRep := repository.NewUserRepository()
	if *middleware.GetUserId(sessionID) == constants.GuestID {
		userRep = repository.NewUserRepository(repository.GuestMode)
	}
	if err != nil {
		return openapi_models.GetUserInfoResponse{}
	}
	userId := middleware.GetUserId(sessionID)

	var userInfoResponse openapi_models.GetUserInfoResponse

	// TODO: DI設計の時に GuestRepositoryを作成すること。このつくりはOpenClosedに違反
	if *userId == constants.GuestID {
		guestId := int32(constants.GuestID)
		userInfoResponse = openapi_models.GetUserInfoResponse{
			User: openapi_models.User{
				Id:               &guestId,
				DepartmentId:     -1,
				LimitOfOperation: 0,
				LastName:         "ゲスト",
				FirstName:        "",
				Password:         "",
				Email:            "",
				CreatedAt:        time.Now(),
				UpdatedAt:        0,
				Role:             constants.RoleGuest,
			},
		}
	} else {
		user := userRep.Find(*userId)
		userInfoResponse = openapi_models.GetUserInfoResponse{
			User: openapi_models.User{
				Id:               user.Id,
				DepartmentId:     user.DepartmentId,
				LimitOfOperation: user.LimitOfOperation,
				LastName:         user.LastName,
				FirstName:        user.FirstName,
				Password:         "",
				Email:            user.Email,
				CreatedAt:        user.CreatedAt,
				UpdatedAt:        user.UpdatedAt,
				Role:             user.Role,
			},
		}
	}
	return userInfoResponse
}
