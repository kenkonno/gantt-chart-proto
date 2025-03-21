package detect_work_outside_employment_periods

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/kenkonno/gantt-chart-proto/backend/api/constants"
	"github.com/kenkonno/gantt-chart-proto/backend/api/middleware"
	"github.com/kenkonno/gantt-chart-proto/backend/api/openapi_models"
	"github.com/kenkonno/gantt-chart-proto/backend/api/utils"
	"github.com/kenkonno/gantt-chart-proto/backend/models/db"
	"github.com/kenkonno/gantt-chart-proto/backend/repository"
	"github.com/samber/lo"
	"net/http"
	"time"
)

func GetDetectWorkOutsideEmploymentPeriodsInvoke(c *gin.Context) (openapi_models.GetDetectWorkOutsideEmploymentPeriodsResponse, error) {

	var request openapi_models.GetDetectWorkOutsideEmploymentPeriodsRequest
	if err := c.ShouldBindQuery(&request); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		panic(err)
	}

	ticketRep := repository.NewTicketRepository(middleware.GetRepositoryMode(c)...)
	tickets := ticketRep.FindByUserIds([]int32{request.UserId}, constants.FacilityStatusEnabled)

	EmploymentStartDate := request.EmploymentStartDate
	var EmploymentEndDate *time.Time = nil
	if !request.EmploymentEndDate.IsZero() {
		EmploymentEndDate = &request.EmploymentEndDate
	}

	outsideTickets := utils.DetectWorkOutsideEmploymentPeriods(tickets, EmploymentStartDate, EmploymentEndDate)

	return openapi_models.GetDetectWorkOutsideEmploymentPeriodsResponse{
		List: lo.Map(outsideTickets, func(item db.Ticket, index int) openapi_models.Ticket {
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
	}, nil

}
