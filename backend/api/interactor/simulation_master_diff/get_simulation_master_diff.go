package simulation_master_diff

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/kenkonno/gantt-chart-proto/backend/api/openapi_models"
	"github.com/kenkonno/gantt-chart-proto/backend/models/db"
	"github.com/kenkonno/gantt-chart-proto/backend/repository"
	"github.com/samber/lo"
	"reflect"
)

func GetSimulationMasterDiffInvoke(c *gin.Context) (openapi_models.GetSimulationMasterDiffResponse, error) {

	processRep := repository.NewProcessRepository()
	simulateProcessRep := repository.NewProcessRepository(repository.SimulationMode)

	departmentRep := repository.NewDepartmentRepository()
	simulateProcessDepartmentRep := repository.NewDepartmentRepository(repository.SimulationMode)

	userRep := repository.NewUserRepository()
	simulateUserRep := repository.NewUserRepository(repository.SimulationMode)

	unitRep := repository.NewUnitRepository()
	simulateUnitRep := repository.NewUnitRepository(repository.SimulationMode)

	holidayRep := repository.NewHolidayRepository()
	simulateHolidayRep := repository.NewHolidayRepository(repository.SimulationMode)

	milestoneRep := repository.NewMilestoneRepository()
	simulateMilestoneRep := repository.NewMilestoneRepository(repository.SimulationMode)

	beforeProcesses, diffProcesses, afterProcesses := getBeforeAfter(processRep.FindAll(), simulateProcessRep.FindAll())
	beforeDepartments, diffDepartments, afterDepartments := getBeforeAfter(departmentRep.FindAll(), simulateProcessDepartmentRep.FindAll())
	beforeUsers, diffUsers, afterUsers := getBeforeAfter(userRep.FindAll(), simulateUserRep.FindAll())
	beforeUnits, diffUnits, afterUnits := getBeforeAfter(unitRep.FindAll(), simulateUnitRep.FindAll())
	beforeHolidays, diffHolidays, afterHolidays := getBeforeAfter(holidayRep.FindAll(), simulateHolidayRep.FindAll())
	beforeMilestones, diffMilestones, afterMilestones := getBeforeAfter(milestoneRep.FindAll(), simulateMilestoneRep.FindAll())

	// 更新が合った者
	return openapi_models.GetSimulationMasterDiffResponse{
		BeforeProcesses:   convertProcesses(beforeProcesses),
		DiffProcesses:     convertProcesses(diffProcesses),
		AfterProcesses:    convertProcesses(afterProcesses),
		BeforeDepartments: convertDepartments(beforeDepartments),
		DiffDepartments:   convertDepartments(diffDepartments),
		AfterDepartments:  convertDepartments(afterDepartments),
		BeforeUsers:       convertUsers(beforeUsers),
		DiffUsers:         convertUsers(diffUsers),
		AfterUsers:        convertUsers(afterUsers),
		BeforeUnits:       convertUnits(beforeUnits),
		DiffUnits:         convertUnits(diffUnits),
		AfterUnits:        convertUnits(afterUnits),
		BeforeHolidays:    convertHolidays(beforeHolidays),
		DiffHolidays:      convertHolidays(diffHolidays),
		AfterHolidays:     convertHolidays(afterHolidays),
		BeforeMilestones:  convertMilestones(beforeMilestones),
		DiffMilestones:    convertMilestones(diffMilestones),
		AfterMilestones:   convertMilestones(afterMilestones),
	}, nil
}

func getBeforeAfter[V comparable](preList []V, postList []V) ([]V, []V, []V) {
	preUpdateList := exceptSlice(preList, postList)
	postUpdateList := exceptSlice(postList, preList)
	beforeDiff, afterDiff := getDiffSlice(intersectSlice(preList, postList))

	var diffList []V
	for index, v := range beforeDiff {
		diffList = append(diffList, v)
		diffList = append(diffList, afterDiff[index])
	}

	return preUpdateList, diffList, postUpdateList

}

// exceptSlice preUpdateListに存在しpostUpdateListに存在しない要素を含む新しいスライス、resultを返します。
func exceptSlice[V any](preUpdateList []V, postUpdateList []V) []V {
	preUpdateMap := make(map[interface{}]bool)
	for _, v := range preUpdateList {
		value := reflect.ValueOf(v)
		id := reflect.Indirect(value).FieldByName("Id").Elem().Int()
		preUpdateMap[id] = true
	}
	for _, v := range postUpdateList {
		value := reflect.ValueOf(v)
		id := reflect.Indirect(value).FieldByName("Id").Elem().Int()
		delete(preUpdateMap, id)
	}
	var result []V
	for _, v := range preUpdateList {
		value := reflect.ValueOf(v)
		id := reflect.Indirect(value).FieldByName("Id").Elem().Int()
		if _, exist := preUpdateMap[id]; exist {
			result = append(result, v)
		}
	}
	return result
}

// intersectSlice preUpdateListとpostUpdateListの両方のスライスに存在する要素を含む新しいスライス、resultを返します。
func intersectSlice[V any](preUpdateList []V, postUpdateList []V) ([]V, []V) {
	preUpdateMap := make(map[interface{}]bool)
	for _, v := range preUpdateList {
		value := reflect.ValueOf(v)
		id := reflect.Indirect(value).FieldByName("Id").Elem().Int()
		preUpdateMap[id] = true
	}
	var preResult []V
	var postResult []V
	for _, v := range postUpdateList {
		value := reflect.ValueOf(v)
		id := reflect.Indirect(value).FieldByName("Id").Elem().Int()
		if _, exist := preUpdateMap[id]; exist {
			postResult = append(postResult, v)
			preResult = append(preResult, preUpdateList[getIndexOfID(preUpdateList, id)])
		}
	}
	return preResult, postResult
}

// getIndexOfID returns the index of the object in list whose ID matches the provided id.
func getIndexOfID[V any](list []V, id interface{}) int {
	for i, v := range list {
		value := reflect.ValueOf(v)
		itemId := reflect.Indirect(value).FieldByName("Id").Elem().Int()
		if itemId == id {
			return i
		}
	}
	return -1 // returns -1 if no matching ID was found
}

func getDiffSlice[V comparable](preUpdateList []V, postUpdateList []V) ([]V, []V) {
	var beforeResult, afterResult []V
	for index, item := range preUpdateList {
		if !reflect.DeepEqual(item, postUpdateList[index]) {
			beforeResult = append(beforeResult, item)
			afterResult = append(afterResult, postUpdateList[index])
			fmt.Println("DIFF")
			fmt.Println(item, postUpdateList[index])
		}
	}
	return beforeResult, afterResult
}

func convertProcesses(list []db.Process) []openapi_models.Process {
	return lo.Map(list, func(item db.Process, index int) openapi_models.Process {
		return openapi_models.Process{
			Id:        item.Id,
			Name:      item.Name,
			Order:     int32(item.Order),
			CreatedAt: item.CreatedAt,
			UpdatedAt: item.UpdatedAt,
			Color:     item.Color,
		}
	})
}
func convertDepartments(list []db.Department) []openapi_models.Department {
	return lo.Map(list, func(item db.Department, index int) openapi_models.Department {
		return openapi_models.Department{
			Id:        item.Id,
			Name:      item.Name,
			Color:     item.Color,
			Order:     int32(item.Order),
			CreatedAt: item.CreatedAt,
			UpdatedAt: item.UpdatedAt,
		}
	})
}
func convertUsers(list []db.User) []openapi_models.User {
	return lo.Map(list, func(item db.User, index int) openapi_models.User {
		return openapi_models.User{
			Id:               item.Id,
			DepartmentId:     item.DepartmentId,
			LimitOfOperation: item.LimitOfOperation,
			LastName:         item.LastName,
			FirstName:        item.FirstName,
			Password:         item.Password,
			Email:            item.Email,
			CreatedAt:        item.CreatedAt,
			UpdatedAt:        item.UpdatedAt,
			Role:             item.Role,
		}
	})
}
func convertUnits(list []db.Unit) []openapi_models.Unit {
	return lo.Map(list, func(item db.Unit, index int) openapi_models.Unit {
		return openapi_models.Unit{
			Id:         item.Id,
			FacilityId: item.FacilityId,
			Name:       item.Name,
			Order:      int32(item.Order),
			CreatedAt:  item.CreatedAt,
			UpdatedAt:  item.UpdatedAt,
		}
	})
}
func convertHolidays(list []db.Holiday) []openapi_models.Holiday {
	return lo.Map(list, func(item db.Holiday, index int) openapi_models.Holiday {
		return openapi_models.Holiday{
			Id:         item.Id,
			Name:       item.Name,
			Date:       item.Date,
			CreatedAt:  item.CreatedAt,
			UpdatedAt:  item.UpdatedAt,
			FacilityId: item.FacilityId,
		}
	})
}
func convertMilestones(list []db.Milestone) []openapi_models.Milestone {
	return lo.Map(list, func(item db.Milestone, index int) openapi_models.Milestone {
		return openapi_models.Milestone{
			Id:          item.Id,
			FacilityId:  item.FacilityId,
			Date:        item.Date,
			Description: item.Description,
			Order:       int32(item.Order),
			CreatedAt:   item.CreatedAt,
			UpdatedAt:   int(item.UpdatedAt),
		}
	})
}
