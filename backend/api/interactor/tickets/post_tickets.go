package tickets

import (
	"github.com/gin-gonic/gin"
	"github.com/kenkonno/gantt-chart-proto/backend/api/openapi_models"
	"github.com/kenkonno/gantt-chart-proto/backend/models/db"
	"github.com/kenkonno/gantt-chart-proto/backend/repository"
	"net/http"
	"time"
)

func PostTicketsInvoke(c *gin.Context) openapi_models.PostTicketsResponse {

	ticketRep := repository.NewTicketRepository()

	var ticketReq openapi_models.PostTicketsRequest
	if err := c.ShouldBindJSON(&ticketReq); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		panic(err)
	}
	result := ticketRep.Upsert(db.Ticket{
		GanttGroupId:    ticketReq.Ticket.GanttGroupId,
		ProcessId:       ticketReq.Ticket.ProcessId,
		DepartmentId:    ticketReq.Ticket.DepartmentId,
		LimitDate:       ticketReq.Ticket.LimitDate,
		Estimate:        ticketReq.Ticket.Estimate,
		DaysAfter:       ticketReq.Ticket.DaysAfter,
		StartDate:       ticketReq.Ticket.StartDate,
		EndDate:         ticketReq.Ticket.EndDate,
		ProgressPercent: ticketReq.Ticket.ProgressPercent,
		Order:           ticketReq.Ticket.Order,
		CreatedAt:       time.Time{},
		UpdatedAt:       0,
	})

	return openapi_models.PostTicketsResponse{
		Ticket: openapi_models.Ticket{
			Id:              result.Id,
			GanttGroupId:    result.GanttGroupId,
			ProcessId:       result.ProcessId,
			DepartmentId:    result.DepartmentId,
			LimitDate:       result.LimitDate,
			Estimate:        result.Estimate,
			DaysAfter:       result.DaysAfter,
			StartDate:       result.StartDate,
			EndDate:         result.EndDate,
			ProgressPercent: result.ProgressPercent,
			Order:           result.Order,
			CreatedAt:       result.CreatedAt,
			UpdatedAt:       result.UpdatedAt,
		},
	}

}
