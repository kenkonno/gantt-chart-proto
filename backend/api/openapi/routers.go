/*
 * GanttChartApi
 *
 * No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)
 *
 * API version: 1.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package openapi

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Route is the information for every URI.
type Route struct {
	// Name is the name of this Route.
	Name string
	// Method is the string for the HTTP method. ex) GET, POST etc..
	Method string
	// Pattern is the pattern of the URI.
	Pattern string
	// HandlerFunc is the handler function of this route.
	HandlerFunc gin.HandlerFunc
}

// Routes is the list of the generated Route.
type Routes []Route

// NewRouter returns a new router.
func NewRouter(router *gin.Engine) *gin.Engine {

	for _, route := range routes {
		switch route.Method {
		case http.MethodGet:
			router.GET(route.Pattern, route.HandlerFunc)
		case http.MethodPost:
			router.POST(route.Pattern, route.HandlerFunc)
		case http.MethodPut:
			router.PUT(route.Pattern, route.HandlerFunc)
		case http.MethodPatch:
			router.PATCH(route.Pattern, route.HandlerFunc)
		case http.MethodDelete:
			router.DELETE(route.Pattern, route.HandlerFunc)
		}
	}

	return router
}

// Index is the index handler.
func Index(c *gin.Context) {
	c.String(http.StatusOK, "Hello World!")
}

var routes = Routes{
	{
		"Index",
		http.MethodGet,
		"/",
		Index,
	},

	{
		"DeleteDepartmentsId",
		http.MethodDelete,
		"/api/departments/:id",
		DeleteDepartmentsId,
	},

	{
		"DeleteFacilitiesId",
		http.MethodDelete,
		"/api/facilities/:id",
		DeleteFacilitiesId,
	},

	{
		"DeleteFacilitySharedLinksId",
		http.MethodDelete,
		"/api/facilitySharedLinks/:id",
		DeleteFacilitySharedLinksId,
	},

	{
		"DeleteFacilityWorkSchedulesId",
		http.MethodDelete,
		"/api/facilityWorkSchedules/:id",
		DeleteFacilityWorkSchedulesId,
	},

	{
		"DeleteGanttGroupsId",
		http.MethodDelete,
		"/api/ganttGroups/:id",
		DeleteGanttGroupsId,
	},

	{
		"DeleteHolidaysId",
		http.MethodDelete,
		"/api/holidays/:id",
		DeleteHolidaysId,
	},

	{
		"DeleteMilestonesId",
		http.MethodDelete,
		"/api/milestones/:id",
		DeleteMilestonesId,
	},

	{
		"DeleteOperationSettingsId",
		http.MethodDelete,
		"/api/operationSettings/:id",
		DeleteOperationSettingsId,
	},

	{
		"DeleteProcessesId",
		http.MethodDelete,
		"/api/processes/:id",
		DeleteProcessesId,
	},

	{
		"DeleteSimulation",
		http.MethodDelete,
		"/api/simulation",
		DeleteSimulation,
	},

	{
		"DeleteTicketUsersId",
		http.MethodDelete,
		"/api/ticketUsers/:id",
		DeleteTicketUsersId,
	},

	{
		"DeleteTicketsId",
		http.MethodDelete,
		"/api/tickets/:id",
		DeleteTicketsId,
	},

	{
		"DeleteUnitsId",
		http.MethodDelete,
		"/api/units/:id",
		DeleteUnitsId,
	},

	{
		"DeleteUsersId",
		http.MethodDelete,
		"/api/users/:id",
		DeleteUsersId,
	},

	{
		"GetAllTickets",
		http.MethodGet,
		"/api/all-tickets",
		GetAllTickets,
	},

	{
		"GetDefaultPileUps",
		http.MethodGet,
		"/api/defaultPileUps",
		GetDefaultPileUps,
	},

	{
		"GetDepartments",
		http.MethodGet,
		"/api/departments",
		GetDepartments,
	},

	{
		"GetDepartmentsId",
		http.MethodGet,
		"/api/departments/:id",
		GetDepartmentsId,
	},

	{
		"GetDetectWorkOutsideEmploymentPeriods",
		http.MethodGet,
		"/api/detectWorkOutsideEmploymentPeriods",
		GetDetectWorkOutsideEmploymentPeriods,
	},

	{
		"GetFacilities",
		http.MethodGet,
		"/api/facilities",
		GetFacilities,
	},

	{
		"GetFacilitiesId",
		http.MethodGet,
		"/api/facilities/:id",
		GetFacilitiesId,
	},

	{
		"GetFacilitySharedLinks",
		http.MethodGet,
		"/api/facilitySharedLinks",
		GetFacilitySharedLinks,
	},

	{
		"GetFacilitySharedLinksId",
		http.MethodGet,
		"/api/facilitySharedLinks/:id",
		GetFacilitySharedLinksId,
	},

	{
		"GetFacilityWorkSchedules",
		http.MethodGet,
		"/api/facilityWorkSchedules",
		GetFacilityWorkSchedules,
	},

	{
		"GetFacilityWorkSchedulesId",
		http.MethodGet,
		"/api/facilityWorkSchedules/:id",
		GetFacilityWorkSchedulesId,
	},

	{
		"GetGanttGroups",
		http.MethodGet,
		"/api/ganttGroups",
		GetGanttGroups,
	},

	{
		"GetGanttGroupsId",
		http.MethodGet,
		"/api/ganttGroups/:id",
		GetGanttGroupsId,
	},

	{
		"GetHolidays",
		http.MethodGet,
		"/api/holidays",
		GetHolidays,
	},

	{
		"GetHolidaysId",
		http.MethodGet,
		"/api/holidays/:id",
		GetHolidaysId,
	},

	{
		"GetMilestones",
		http.MethodGet,
		"/api/milestones",
		GetMilestones,
	},

	{
		"GetMilestonesId",
		http.MethodGet,
		"/api/milestones/:id",
		GetMilestonesId,
	},

	{
		"GetOperationSettingsId",
		http.MethodGet,
		"/api/operationSettings/:id",
		GetOperationSettingsId,
	},

	{
		"GetPileUps",
		http.MethodGet,
		"/api/pileUps",
		GetPileUps,
	},

	{
		"GetProcesses",
		http.MethodGet,
		"/api/processes",
		GetProcesses,
	},

	{
		"GetProcessesId",
		http.MethodGet,
		"/api/processes/:id",
		GetProcessesId,
	},

	{
		"GetScheduleAlerts",
		http.MethodGet,
		"/api/scheduleAlerts",
		GetScheduleAlerts,
	},

	{
		"GetSimulation",
		http.MethodGet,
		"/api/simulation",
		GetSimulation,
	},

	{
		"GetSimulationMasterDiff",
		http.MethodGet,
		"/api/simulationMasterDiff",
		GetSimulationMasterDiff,
	},

	{
		"GetTicketMemoId",
		http.MethodGet,
		"/api/ticket-memo/:id",
		GetTicketMemoId,
	},

	{
		"GetTicketUsers",
		http.MethodGet,
		"/api/ticketUsers",
		GetTicketUsers,
	},

	{
		"GetTicketUsersId",
		http.MethodGet,
		"/api/ticketUsers/:id",
		GetTicketUsersId,
	},

	{
		"GetTickets",
		http.MethodGet,
		"/api/tickets",
		GetTickets,
	},

	{
		"GetTicketsId",
		http.MethodGet,
		"/api/tickets/:id",
		GetTicketsId,
	},

	{
		"GetUnits",
		http.MethodGet,
		"/api/units",
		GetUnits,
	},

	{
		"GetUnitsId",
		http.MethodGet,
		"/api/units/:id",
		GetUnitsId,
	},

	{
		"GetUserInfo",
		http.MethodGet,
		"/api/userInfo",
		GetUserInfo,
	},

	{
		"GetUsers",
		http.MethodGet,
		"/api/users",
		GetUsers,
	},

	{
		"GetUsersId",
		http.MethodGet,
		"/api/users/:id",
		GetUsersId,
	},

	{
		"PostBulkUpdateTickets",
		http.MethodPost,
		"/api/bulkUpdateTickets",
		PostBulkUpdateTickets,
	},

	{
		"PostCopyFacilitys",
		http.MethodPost,
		"/api/copyFacilitys",
		PostCopyFacilitys,
	},

	{
		"PostDepartments",
		http.MethodPost,
		"/api/departments",
		PostDepartments,
	},

	{
		"PostDepartmentsId",
		http.MethodPost,
		"/api/departments/:id",
		PostDepartmentsId,
	},

	{
		"PostFacilities",
		http.MethodPost,
		"/api/facilities",
		PostFacilities,
	},

	{
		"PostFacilitiesId",
		http.MethodPost,
		"/api/facilities/:id",
		PostFacilitiesId,
	},

	{
		"PostFacilitySharedLinks",
		http.MethodPost,
		"/api/facilitySharedLinks",
		PostFacilitySharedLinks,
	},

	{
		"PostFacilitySharedLinksId",
		http.MethodPost,
		"/api/facilitySharedLinks/:id",
		PostFacilitySharedLinksId,
	},

	{
		"PostFacilityWorkSchedules",
		http.MethodPost,
		"/api/facilityWorkSchedules",
		PostFacilityWorkSchedules,
	},

	{
		"PostFacilityWorkSchedulesId",
		http.MethodPost,
		"/api/facilityWorkSchedules/:id",
		PostFacilityWorkSchedulesId,
	},

	{
		"PostGanttGroups",
		http.MethodPost,
		"/api/ganttGroups",
		PostGanttGroups,
	},

	{
		"PostGanttGroupsId",
		http.MethodPost,
		"/api/ganttGroups/:id",
		PostGanttGroupsId,
	},

	{
		"PostHolidays",
		http.MethodPost,
		"/api/holidays",
		PostHolidays,
	},

	{
		"PostHolidaysId",
		http.MethodPost,
		"/api/holidays/:id",
		PostHolidaysId,
	},

	{
		"PostLogin",
		http.MethodPost,
		"/api/login",
		PostLogin,
	},

	{
		"PostLogout",
		http.MethodPost,
		"/api/logout",
		PostLogout,
	},

	{
		"PostMilestones",
		http.MethodPost,
		"/api/milestones",
		PostMilestones,
	},

	{
		"PostMilestonesId",
		http.MethodPost,
		"/api/milestones/:id",
		PostMilestonesId,
	},

	{
		"PostOperationSettingsId",
		http.MethodPost,
		"/api/operationSettings/:id",
		PostOperationSettingsId,
	},

	{
		"PostProcesses",
		http.MethodPost,
		"/api/processes",
		PostProcesses,
	},

	{
		"PostProcessesId",
		http.MethodPost,
		"/api/processes/:id",
		PostProcessesId,
	},

	{
		"PostSimulation",
		http.MethodPost,
		"/api/simulation",
		PostSimulation,
	},

	{
		"PostTicketMemoId",
		http.MethodPost,
		"/api/ticket-memo/:id",
		PostTicketMemoId,
	},

	{
		"PostTicketUsers",
		http.MethodPost,
		"/api/ticketUsers",
		PostTicketUsers,
	},

	{
		"PostTickets",
		http.MethodPost,
		"/api/tickets",
		PostTickets,
	},

	{
		"PostTicketsId",
		http.MethodPost,
		"/api/tickets/:id",
		PostTicketsId,
	},

	{
		"PostUnits",
		http.MethodPost,
		"/api/units",
		PostUnits,
	},

	{
		"PostUnitsDuplicate",
		http.MethodPost,
		"/api/units/duplicate",
		PostUnitsDuplicate,
	},

	{
		"PostUnitsId",
		http.MethodPost,
		"/api/units/:id",
		PostUnitsId,
	},

	{
		"PostUploadUsersCsvFile",
		http.MethodPost,
		"/api/users/upload",
		PostUploadUsersCsvFile,
	},

	{
		"PostUsers",
		http.MethodPost,
		"/api/users",
		PostUsers,
	},

	{
		"PostUsersId",
		http.MethodPost,
		"/api/users/:id",
		PostUsersId,
	},

	{
		"PutSimulation",
		http.MethodPut,
		"/api/simulation",
		PutSimulation,
	},
}
