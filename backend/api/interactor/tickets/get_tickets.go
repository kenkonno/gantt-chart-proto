package tickets

import (
	"github.com/gin-gonic/gin"
	"github.com/kenkonno/gantt-chart-proto/backend/api/middleware"
	"github.com/kenkonno/gantt-chart-proto/backend/api/openapi_models"
	"github.com/kenkonno/gantt-chart-proto/backend/models/db"
	"github.com/kenkonno/gantt-chart-proto/backend/repository"
	"github.com/samber/lo"
	"strconv"
)

func GetTicketsInvoke(c *gin.Context) (openapi_models.GetTicketsResponse, error) {
	ganttGroupIdsParam := c.QueryArray("ganttGroupIds")
	// シミュレーション中に本番のチケットを取得するために用意
	mode := c.Query("mode")

	ticketRep := repository.NewTicketRepository(middleware.GetRepositoryMode(c)...)
	if mode == "prod" {
		ticketRep = repository.NewTicketRepository()
	}

	var ganttGroupIds []int32
	for _, v := range ganttGroupIdsParam {
		vv, err := strconv.Atoi(v)
		if err != nil {
			panic(err)
		}
		ganttGroupIds = append(ganttGroupIds, int32(vv))
	}

	ticketList := ticketRep.FindByGanttGroupIds(ganttGroupIds)

	return openapi_models.GetTicketsResponse{
		List: lo.Map(ticketList, func(item db.Ticket, index int) openapi_models.Ticket {
			// ゲストの場合はマスクをする
			if middleware.IsGuest(c) {
				return openapi_models.Ticket{
					Id:              item.Id,
					GanttGroupId:    item.GanttGroupId,
					ProcessId:       item.ProcessId,
					DepartmentId:    item.DepartmentId,
					LimitDate:       nil,
					Estimate:        nil,
					NumberOfWorker:  nil,
					DaysAfter:       nil,
					StartDate:       item.StartDate,
					EndDate:         item.EndDate,
					ProgressPercent: nil,
					Order:           item.Order,
					CreatedAt:       item.CreatedAt,
					UpdatedAt:       item.UpdatedAt,
				}
			} else {
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
			}
		}),
	}, nil
}
