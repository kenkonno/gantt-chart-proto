package user
import (
	"github.com/gin-gonic/gin"
	"github.com/kenkonno/gantt-chart-proto/backend/api/openapi_models"
	"github.com/kenkonno/gantt-chart-proto/backend/repository"
)
func GetUsersInvoke(c *gin.Context) openapi_models.GetUsersResponse {
	userRep := repository.NewUserRepository()

	userList := userRep.FindAll()

	return openapi_models.GetUsersResponse{}
}
