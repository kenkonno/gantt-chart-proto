package schedule_alerts

import (
	"github.com/gin-gonic/gin"
	"github.com/kenkonno/gantt-chart-proto/backend/api/middleware"
	"github.com/kenkonno/gantt-chart-proto/backend/api/openapi_models"
	"github.com/kenkonno/gantt-chart-proto/backend/models/db"
	"github.com/kenkonno/gantt-chart-proto/backend/repository"
	"github.com/samber/lo"
)

func GetScheduleAlertsInvoke(c *gin.Context) openapi_models.GetScheduleAlertsResponse {
	scheduleAlertRep := repository.NewScheduleAlertRepository(middleware.GetRepositoryMode(c)...)
	scheduleAlerts := scheduleAlertRep.FindAll()
	return openapi_models.GetScheduleAlertsResponse{
		List: lo.Map(scheduleAlerts, func(item db.ScheduleAlert, index int) openapi_models.ScheduleAlert {
			return openapi_models.ScheduleAlert{
				FacilityId:         item.FacilityId,
				FacilityName:       item.FacilityName,
				UnitId:             item.UnitId,
				UnitName:           item.UnitName,
				ProcessId:          item.ProcessId,
				ProcessName:        item.ProcessName,
				TicketId:           item.TicketId,
				EndDate:            item.EndDate,
				StartDate:          item.StartDate,
				ActualProgressDate: item.ActualProgressDate, // TODO: 未着手に適当なものが入っている
				ProgressPercent:    item.ProgressPercent,
				DelayDays:          item.DelayDays,
			}
		}),
	}
}
