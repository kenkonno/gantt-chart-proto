package tickets

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/kenkonno/gantt-chart-proto/backend/api/openapi_models"
	"github.com/kenkonno/gantt-chart-proto/backend/models/db"
	"github.com/kenkonno/gantt-chart-proto/backend/repository"
	"github.com/kenkonno/gantt-chart-proto/backend/repository/connection"
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

	result, err := ticketRep.Upsert(db.Ticket{
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
		UpdatedAt:       ticketReq.Ticket.UpdatedAt,
	})
	if err != nil {
		var target connection.ConflictError
		if errors.As(err, &target) {
			c.JSON(http.StatusConflict, err.Error())
			panic(err)
		}
	}

	return openapi_models.PostTicketsIdResponse{
		Ticket: openapi_models.Ticket{
			Id:              result.Id,
			GanttGroupId:    result.GanttGroupId,
			ProcessId:       result.ProcessId,
			DepartmentId:    result.DepartmentId,
			LimitDate:       result.LimitDate,
			Estimate:        result.Estimate,
			NumberOfWorker:  result.NumberOfWorker,
			DaysAfter:       result.DaysAfter,
			StartDate:       result.StartDate,
			EndDate:         result.EndDate,
			ProgressPercent: result.ProgressPercent,
			Memo:            result.Memo,
			Order:           result.Order,
			CreatedAt:       result.CreatedAt,
			UpdatedAt:       result.UpdatedAt,
		},
	}

}
