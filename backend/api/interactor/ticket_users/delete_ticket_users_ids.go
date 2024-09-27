package ticket_users

import (
	"github.com/gin-gonic/gin"
	"github.com/kenkonno/gantt-chart-proto/backend/api/middleware"
	"github.com/kenkonno/gantt-chart-proto/backend/api/openapi_models"
	"github.com/kenkonno/gantt-chart-proto/backend/repository"
	"strconv"
)

func DeleteTicketUsersIdInvoke(c *gin.Context) openapi_models.DeleteTicketUsersIdResponse {

	ticketUserRep := repository.NewTicketUserRepository(middleware.GetRepositoryMode(c)...)

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		panic(err)
	}

	ticketUserRep.Delete(int32(id))

	return openapi_models.DeleteTicketUsersIdResponse{}

}
