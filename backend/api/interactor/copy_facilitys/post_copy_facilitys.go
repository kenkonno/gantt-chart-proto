package copy_facilitys

import (
	"github.com/gin-gonic/gin"
	"github.com/kenkonno/gantt-chart-proto/backend/api/openapi_models"
	"github.com/kenkonno/gantt-chart-proto/backend/models/db"
	"github.com/kenkonno/gantt-chart-proto/backend/repository"
	"github.com/samber/lo"
	"net/http"
	"time"
)

func PostCopyFacilitysInvoke(c *gin.Context) openapi_models.PostCopyFacilitysResponse {

	// 設備のコピーを実施する
	facilityRep := repository.NewFacilityRepository()
	ganttGroupRep := repository.NewGanttGroupRepository()
	unitRep := repository.NewUnitRepository()
	ticketRep := repository.NewTicketRepository()

	var copyFacilityReq openapi_models.PostCopyFacilitysRequest
	if err := c.ShouldBindJSON(&copyFacilityReq); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		panic(err)
	}

	// コピー元のFacility
	orgFacility := facilityRep.Find(copyFacilityReq.FacilityId)
	// コピー元のGanttGroups
	orgGanttGroups := ganttGroupRep.FindByFacilityId(*orgFacility.Id)
	// コピー元のunit
	orgUnits := unitRep.FindByFacilityId(*orgFacility.Id)
	// コピー元のTicket
	orgTickets := ticketRep.FindByGanttGroupIds(lo.Map(orgGanttGroups, func(item db.GanttGroup, index int) int32 {
		return item.FacilityId
	}))

	// コピーを実施する、順番は Facility -> Unit | GanttGroups -> (Ticket)
	newFacility := facilityRep.Upsert(db.Facility{
		Id:        nil,
		Name:      copyFacilityReq.Facility.Name,
		TermFrom:  copyFacilityReq.Facility.TermFrom,
		TermTo:    copyFacilityReq.Facility.TermTo,
		Order:     int(copyFacilityReq.Facility.Order),
		CreatedAt: time.Time{},
		UpdatedAt: 0,
	})

	var unitMap = make(map[int32]int32)
	var ganttGroupMap = make(map[int32]int32)

	// ユニットのコピー
	for _, unit := range orgUnits {
		newUnit := unitRep.Upsert(db.Unit{
			Id:         nil,
			Name:       unit.Name,
			FacilityId: *newFacility.Id,
			Order:      unit.Order,
			CreatedAt:  time.Time{},
			UpdatedAt:  0,
		})
		unitMap[*unit.Id] = *newUnit.Id
	}

	// ganttGroupのコピー
	for _, group := range orgGanttGroups {
		newGanttGroup := ganttGroupRep.Upsert(db.GanttGroup{
			Id:         nil,
			FacilityId: *newFacility.Id,
			UnitId:     unitMap[group.UnitId],
			CreatedAt:  time.Time{},
			UpdatedAt:  0,
		})
		ganttGroupMap[*group.Id] = *newGanttGroup.Id
	}

	// Ticketのコピー
	for _, ticket := range orgTickets {
		ticketRep.Upsert(db.Ticket{
			Id:              nil,
			GanttGroupId:    ganttGroupMap[ticket.GanttGroupId],
			ProcessId:       ticket.ProcessId,
			DepartmentId:    ticket.DepartmentId,
			LimitDate:       nil,
			Estimate:        nil,
			NumberOfWorker:  nil,
			DaysAfter:       nil,
			StartDate:       &newFacility.TermFrom,
			EndDate:         &newFacility.TermTo,
			ProgressPercent: nil,
			Order:           ticket.Order,
			CreatedAt:       time.Time{},
			UpdatedAt:       0,
		})
	}

	return openapi_models.PostCopyFacilitysResponse{}

}
