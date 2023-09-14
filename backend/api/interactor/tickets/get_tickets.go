package tickets

import (
	"github.com/gin-gonic/gin"
	"github.com/kenkonno/gantt-chart-proto/backend/api/openapi_models"
	"github.com/kenkonno/gantt-chart-proto/backend/models/db"
	"github.com/kenkonno/gantt-chart-proto/backend/repository"
	"github.com/samber/lo"
	"strconv"
)

func GetTicketsInvoke(c *gin.Context) openapi_models.GetTicketsResponse {
	ticketRep := repository.NewTicketRepository()

	// TODO: メモ疲れたのでもうやめ。不要なAPIは有りそうなので精査する。unit追加時に gantt_groupsも追加するようにした。（そもそもこれもイランかもしれんけど・・
	// TODO: 画面からgantt_groupsと tickets, units のAPIコールして描画するところまで頑張ってやってください。

	// TODO: GETで配列ってどうする？
	ganttGroupIdsParam := c.QueryArray("ganttGroupIds")
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
