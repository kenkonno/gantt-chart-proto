package users

import (
	"github.com/gin-gonic/gin"
	"github.com/kenkonno/gantt-chart-proto/backend/api/openapi_models"
	"github.com/kenkonno/gantt-chart-proto/backend/repository"
	"strconv"
)

func DeleteUsersIdInvoke(c *gin.Context) openapi_models.DeleteUsersIdResponse {

	userRep := repository.NewUserRepository()
	ticketUserRep := repository.NewTicketUserRepository()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		panic(err)
	}

	userRep.Delete(int32(id))

	allTicketUser := ticketUserRep.FindAll()
	for _, item := range allTicketUser {
		if item.UserId == int32(id) {
			ticketUserRep.Delete(*item.Id)
		}
	}

	return openapi_models.DeleteUsersIdResponse{}

}
