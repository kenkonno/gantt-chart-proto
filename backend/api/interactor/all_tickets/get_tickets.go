package all_tickets

import (
	"github.com/gin-gonic/gin"
	"github.com/kenkonno/gantt-chart-proto/backend/api/constants"
	"github.com/kenkonno/gantt-chart-proto/backend/api/openapi_models"
	"github.com/kenkonno/gantt-chart-proto/backend/models/db"
	"github.com/kenkonno/gantt-chart-proto/backend/repository"
	"github.com/samber/lo"
	"golang.org/x/exp/slices"
)

func GetAllTicketsInvoke(c *gin.Context) openapi_models.GetAllTicketsResponse {

	qFacilityTypes := c.QueryArray("facilityTypes")

	var facilityTypes []string

	if slices.Contains(qFacilityTypes, constants.FacilityTypeOrdered) {
		facilityTypes = append(facilityTypes, constants.FacilityTypeOrdered)
	}
	if slices.Contains(qFacilityTypes, constants.FacilityTypePrepared) {
		facilityTypes = append(facilityTypes, constants.FacilityTypePrepared)
	}
	ticketRep := repository.NewTicketRepository()
	ticketList := ticketRep.FindByFacilityType(facilityTypes, []string{constants.FacilityStatusEnabled})

	return openapi_models.GetAllTicketsResponse{
		List: lo.Map(ticketList, func(item db.Ticket, index int) openapi_models.Ticket {
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
				Order:           item.Order,
				CreatedAt:       item.CreatedAt,
				UpdatedAt:       item.UpdatedAt,
			}
		}),
	}
}
