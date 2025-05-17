package units_duplicate

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/kenkonno/gantt-chart-proto/backend/api/middleware"
	"github.com/kenkonno/gantt-chart-proto/backend/api/openapi_models"
	"github.com/kenkonno/gantt-chart-proto/backend/models/db"
	"github.com/kenkonno/gantt-chart-proto/backend/repository"
	"github.com/samber/lo"
	"net/http"
	"strings"
	"time"
)

func PostUnitsDuplicateInvoke(c *gin.Context) (openapi_models.PostUnitsDuplicateResponse, error) {

	unitRep := repository.NewUnitRepository(middleware.GetRepositoryMode(c)...)
	ticketRep := repository.NewTicketRepository(middleware.GetRepositoryMode(c)...)
	ganttGroupsRep := repository.NewGanttGroupRepository(middleware.GetRepositoryMode(c)...)
	ticketUserRep := repository.NewTicketUserRepository(middleware.GetRepositoryMode(c)...)

	var unitReq openapi_models.PostUnitsDuplicateRequest
	if err := c.ShouldBindJSON(&unitReq); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		panic(err)
	}
	unit := unitRep.Find(unitReq.UnitId)

	unit.Id = nil
	unit.Name = unitReq.UnitName

	ganttGroups := ganttGroupsRep.FindByFacilityId([]int32{unit.FacilityId})
	targetGanttGroup, exists := lo.Find(ganttGroups, func(item db.GanttGroup) bool {
		return item.UnitId == unitReq.UnitId
	})
	if !exists {
		return openapi_models.PostUnitsDuplicateResponse{}, errors.New("不正なユニットです。")
	}

	// ユニット紐づくチケットの複製
	tickets := ticketRep.FindByGanttGroupIds([]int32{*targetGanttGroup.Id})

	// ユニットの追加
	r := unitRep.Upsert(db.Unit{
		Name:       strings.TrimSpace(unit.Name),
		FacilityId: unit.FacilityId,
		Order:      unit.Order + 1,
		CreatedAt:  time.Time{},
		UpdatedAt:  0,
	})

	// ganttGroupsの追加
	g := ganttGroupsRep.Upsert(db.GanttGroup{
		Id:         nil,
		FacilityId: unit.FacilityId,
		UnitId:     *r.Id,
		CreatedAt:  time.Time{},
		UpdatedAt:  0,
	})

	for _, ticket := range tickets {
		orgTicketId := *ticket.Id
		ticket.Id = nil
		ticket.GanttGroupId = *g.Id
		// リクエストパラメータのcopyXXXフラグに基づいて値を設定
		if !unitReq.CopyProcessId {
			ticket.ProcessId = nil
		}
		if !unitReq.CopyDepartmentId {
			ticket.DepartmentId = nil
		}
		if !unitReq.CopyEstimate {
			ticket.Estimate = nil
		}
		if !unitReq.CopyNumberOfWorker {
			*ticket.NumberOfWorker = 1
		}
		if !unitReq.CopyDaysAfter {
			ticket.DaysAfter = nil
		}
		if !unitReq.CopyStartDate {
			ticket.StartDate = nil
		}
		if !unitReq.CopyEndDate {
			ticket.EndDate = nil
		}
		if !unitReq.CopyProgressPercent {
			ticket.ProgressPercent = nil
		}
		if !unitReq.CopyMemo {
			ticket.Memo = ""
		}
		newTicket, _ := ticketRep.Upsert(ticket)

		if unitReq.CopyTicketUser {
			ticketUsers := ticketUserRep.FindByTicketId(orgTicketId)
			for _, ticketUser := range ticketUsers {
				ticketUser.Id = nil
				ticketUser.TicketId = *newTicket.Id
				ticketUserRep.Upsert(ticketUser)
			}
		}
	}

	return openapi_models.PostUnitsDuplicateResponse{}, nil

}
