package facilities

import (
	"github.com/gin-gonic/gin"
	"github.com/kenkonno/gantt-chart-proto/backend/api/middleware"
	"github.com/kenkonno/gantt-chart-proto/backend/api/openapi_models"
	"github.com/kenkonno/gantt-chart-proto/backend/models/db"
	"github.com/kenkonno/gantt-chart-proto/backend/repository"
	"net/http"
	"strings"
	"time"
)

func PostFacilitiesInvoke(c *gin.Context) (openapi_models.PostFacilitiesResponse, error) {

	facilityRep := repository.NewFacilityRepository(middleware.GetRepositoryMode(c)...)
	unitRep := repository.NewUnitRepository(middleware.GetRepositoryMode(c)...)
	ganttGroupsRep := repository.NewGanttGroupRepository(middleware.GetRepositoryMode(c)...)
	processRep := repository.NewProcessRepository(middleware.GetRepositoryMode(c)...)
	ticketRep := repository.NewTicketRepository(middleware.GetRepositoryMode(c)...)

	var facilityReq openapi_models.PostFacilitiesRequest
	if err := c.ShouldBindJSON(&facilityReq); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		panic(err)
	}
	newFacility := facilityRep.Upsert(db.Facility{
		Name:            strings.TrimSpace(facilityReq.Facility.Name),
		TermFrom:        facilityReq.Facility.TermFrom,
		TermTo:          facilityReq.Facility.TermTo,
		Order:           int(facilityReq.Facility.Order),
		Status:          facilityReq.Facility.Status,
		Type:            facilityReq.Facility.Type,
		FreeText: facilityReq.Facility.FreeText,
		ShipmentDueDate: facilityReq.Facility.ShipmentDueDate,
		CreatedAt:       time.Time{},
		UpdatedAt:       0,
	})
	// Holiday creation is now facility-independent
	// TODO: post_units.goと重複コード 本体ユニットをデフォルトで登録する
	r := unitRep.Upsert(db.Unit{
		Name:       "本体",
		FacilityId: *newFacility.Id,
		Order:      1,
		CreatedAt:  time.Time{},
		UpdatedAt:  0,
	})
	newGanttGroup := ganttGroupsRep.Upsert(db.GanttGroup{
		Id:         nil,
		FacilityId: *newFacility.Id,
		UnitId:     *r.Id,
		CreatedAt:  time.Time{},
		UpdatedAt:  0,
	})
	// 全行程を登録して、お尻から登録する。開始日終了日は工程の数の案件の期間で決定する
	// 1週刊づつずらして、開始日はMAX(開始日, 案件開始日）とする。
	allProcess := processRep.FindAll()
	for i, v := range allProcess {
		index := len(allProcess) - 1 - i
		endDate := newFacility.TermTo.Add(-1 * (time.Hour * 24) * (7 * time.Duration(index)))
		startDate := endDate.Add(-1 * (time.Hour * 24) * 6)
		if newFacility.TermFrom.After(endDate) {
			endDate = newFacility.TermFrom
		}
		if newFacility.TermFrom.After(startDate) {
			startDate = newFacility.TermFrom
		}
		nOfW := int32(1)                   // 人数はデフォルト1
		_, _ = ticketRep.Upsert(db.Ticket{ // 新規なのでエラーハンドリング無
			GanttGroupId:    *newGanttGroup.Id,
			ProcessId:       v.Id,
			DepartmentId:    nil,
			LimitDate:       nil,
			Estimate:        nil,
			NumberOfWorker:  &nOfW,
			DaysAfter:       nil,
			StartDate:       &startDate,
			EndDate:         &endDate,
			ProgressPercent: nil,
			Order:           int32(i),
			CreatedAt:       time.Time{},
			UpdatedAt:       0,
		})
	}
	return openapi_models.PostFacilitiesResponse{}, nil

}
