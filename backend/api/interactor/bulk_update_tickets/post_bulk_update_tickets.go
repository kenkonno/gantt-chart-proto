package bulk_update_tickets

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/kenkonno/gantt-chart-proto/backend/api/openapi_models"
	"github.com/kenkonno/gantt-chart-proto/backend/models/db"
	"github.com/kenkonno/gantt-chart-proto/backend/repository"
	"github.com/kenkonno/gantt-chart-proto/backend/repository/connection"
	"github.com/samber/lo"
	"net/http"
)

func PostBulkUpdateTicketsInvoke(c *gin.Context) (openapi_models.PostBulkUpdateTicketsResponse, error) {

	ticketRep := repository.NewTicketRepository()

	var bulkUpdateTicketReq openapi_models.PostBulkUpdateTicketsRequest
	if err := c.ShouldBindJSON(&bulkUpdateTicketReq); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		panic(err)
	}

	var results []db.Ticket
	for _, v := range bulkUpdateTicketReq.Tickets {
		ticket := db.Ticket{
			Id:              v.Id,
			GanttGroupId:    v.GanttGroupId,
			ProcessId:       v.ProcessId,
			DepartmentId:    v.DepartmentId,
			LimitDate:       v.LimitDate,
			Estimate:        v.Estimate,
			NumberOfWorker:  v.NumberOfWorker,
			DaysAfter:       v.DaysAfter,
			StartDate:       v.StartDate,
			EndDate:         v.EndDate,
			ProgressPercent: v.ProgressPercent,
			Memo:            v.Memo,
			Order:           v.Order,
			UpdatedAt:       v.UpdatedAt,
		}
		result, err := ticketRep.Upsert(ticket)
		if err != nil {
			var target connection.ConflictError
			if errors.As(err, &target) {
				c.JSON(http.StatusConflict, err.Error())
				panic(err)
			}
		}
		results = append(results, result)
	}

	return openapi_models.PostBulkUpdateTicketsResponse{
		Tickets: lo.Map(results, func(item db.Ticket, index int) openapi_models.Ticket {
			return openapi_models.Ticket{
				Id:              item.Id,
				GanttGroupId:    item.GanttGroupId,
				ProcessId:       item.ProcessId,
				DepartmentId:    item.DepartmentId,
				LimitDate:       item.LimitDate,
				Estimate:        item.Estimate,
				NumberOfWorker:  item.NumberOfWorker,
				DaysAfter:       item.DaysAfter,
				StartDate:       item.StartDate,
				EndDate:         item.EndDate,
				ProgressPercent: item.ProgressPercent,
				Memo:            item.Memo,
				Order:           item.Order,
				CreatedAt:       item.CreatedAt,
				UpdatedAt:       item.UpdatedAt,
			}
		}),
	}, nil

}
