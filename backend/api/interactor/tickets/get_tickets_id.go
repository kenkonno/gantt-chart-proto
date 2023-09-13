package tickets

import (
	"github.com/gin-gonic/gin"
	"github.com/kenkonno/gantt-chart-proto/backend/api/openapi_models"
	"github.com/kenkonno/gantt-chart-proto/backend/repository"
	"strconv"
)

func GetTicketsIdInvoke(c *gin.Context) openapi_models.GetTicketsIdResponse {
	ticketRep := repository.NewTicketRepository()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		panic(err)
	}

	ticket := ticketRep.Find(int32(id))

	return openapi_models.GetTicketsIdResponse{
		Ticket: openapi_models.Ticket{
			Id:              ticket.Id,
			GanttGroupId:    ticket.GanttGroupId,
			ProcessId:       ticket.ProcessId,
			DepartmentId:    ticket.DepartmentId,
			LimitDate:       ticket.LimitDate,
			Estimate:        ticket.Estimate,
			NumberOfWorker:  ticket.NumberOfWorker,
			DaysAfter:       ticket.DaysAfter,
			StartDate:       ticket.StartDate,
			EndDate:         ticket.EndDate,
			ProgressPercent: ticket.ProgressPercent,
			Order:           ticket.Order,
			CreatedAt:       ticket.CreatedAt,
			UpdatedAt:       ticket.UpdatedAt,
		},
	}
}
