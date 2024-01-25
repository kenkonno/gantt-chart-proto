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
	holidayRep := repository.NewHolidayRepository()
	operationSettingRep := repository.NewOperationSettingRepository()

	var copyFacilityReq openapi_models.PostCopyFacilitysRequest
	if err := c.ShouldBindJSON(&copyFacilityReq); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		panic(err)
	}

	// コピー元のFacility
	orgFacility := facilityRep.Find(copyFacilityReq.FacilityId)
	// コピー元のGanttGroups
	orgGanttGroups := ganttGroupRep.FindByFacilityId(*orgFacility.Id)
	// コピー元の稼働設定
	orgOperationSettings := operationSettingRep.FindByFacilityId(*orgFacility.Id)

	// コピー元のunit
	orgUnits := unitRep.FindByFacilityId(*orgFacility.Id)
	// コピー元のTicket
	orgTickets := ticketRep.FindByGanttGroupIds(lo.Map(orgGanttGroups, func(item db.GanttGroup, index int) int32 {
		return *item.Id
	}))

	// コピーを実施する、順番は Facility -> Unit | GanttGroups -> (Ticket)
	allFacility := facilityRep.FindAll([]string{}, []string{})
	order := lo.MaxBy(allFacility,func(a db.Facility, b db.Facility) bool {
		return a.Order > b.Order
	}).Order + 1
	newFacility := facilityRep.Upsert(db.Facility{
		Id:              nil,
		Name:            copyFacilityReq.Facility.Name,
		TermFrom:        copyFacilityReq.Facility.TermFrom,
		TermTo:          copyFacilityReq.Facility.TermTo,
		Order:           order,
		Status:          copyFacilityReq.Facility.Status,
		Type:            copyFacilityReq.Facility.Type,
		ShipmentDueDate: copyFacilityReq.Facility.ShipmentDueDate,
		CreatedAt:       time.Time{},
		UpdatedAt:       0,
	})

	var unitMap = make(map[int32]int32)
	var ganttGroupMap = make(map[int32]int32)

	// ユニットのコピー
	for _, unit := range orgUnits {
		newUnit := unitRep.Upsert(db.Unit{
			Id:         nil,
			Name:       unit.Name + "のコピー",
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

	// 稼働設定のコピー
	for _, setting := range orgOperationSettings {
		operationSettingRep.Upsert(db.OperationSetting{
			Id:         nil,
			FacilityId: *newFacility.Id,
			UnitId:     unitMap[setting.UnitId],
			ProcessId:  setting.ProcessId,
			WorkHour:   setting.WorkHour,
			CreatedAt:  time.Time{},
			UpdatedAt:  0,
		})
	}

	// Ticketのコピー
	currentFrom := newFacility.TermFrom
	for _, ticket := range orgTickets {
		from := currentFrom
		to := from.Add(time.Hour * 24 * 7)
		nOfw := int32(1)
		ticketRep.Upsert(db.Ticket{
			Id:              nil,
			GanttGroupId:    ganttGroupMap[ticket.GanttGroupId],
			ProcessId:       ticket.ProcessId,
			DepartmentId:    ticket.DepartmentId,
			LimitDate:       nil,
			Estimate:        nil,
			NumberOfWorker:  &nOfw,
			DaysAfter:       nil,
			StartDate:       &from,
			EndDate:         &to,
			ProgressPercent: nil,
			Order:           ticket.Order,
			CreatedAt:       time.Time{},
			UpdatedAt:       0,
			
		})
		currentFrom = from.Add(time.Hour * 24)
	}

	// 祝日の投入
	holidayRep.InsertByFacilityId(*newFacility.Id)

	return openapi_models.PostCopyFacilitysResponse{}

}
