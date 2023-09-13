package tickets

import (
	"github.com/gin-gonic/gin"
	"github.com/kenkonno/gantt-chart-proto/backend/api/openapi_models"
	"github.com/kenkonno/gantt-chart-proto/backend/models/db"
	"github.com/kenkonno/gantt-chart-proto/backend/repository"
	"net/http"
	"time"
)

func PostTicketsIdInvoke(c *gin.Context) openapi_models.PostTicketsIdResponse {

	ticketRep := repository.NewTicketRepository()

	var ticketReq openapi_models.PostTicketsRequest
	if err := c.ShouldBindJSON(&ticketReq); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		panic(err)
	}

	ticketRep.Upsert(db.Ticket{
		Id:              ticketReq.Ticket.Id,
		GanttGroupId:    ticketReq.Ticket.GanttGroupId,
		ProcessId:       ticketReq.Ticket.ProcessId,
		DepartmentId:    ticketReq.Ticket.DepartmentId,
		LimitDate:       ticketReq.Ticket.LimitDate,
		Estimate:        ticketReq.Ticket.Estimate,
		NumberOfWorker:  ticketReq.Ticket.NumberOfWorker,
		DaysAfter:       ticketReq.Ticket.DaysAfter,
		StartDate:       ticketReq.Ticket.StartDate,
		EndDate:         ticketReq.Ticket.EndDate,
		ProgressPercent: ticketReq.Ticket.ProgressPercent,
		Order:           ticketReq.Ticket.Order,
		CreatedAt:       time.Time{},
		UpdatedAt:       0,
	})

	return openapi_models.PostTicketsIdResponse{}

}
