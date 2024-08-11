package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/kenkonno/gantt-chart-proto/backend/api/constants"
	"github.com/kenkonno/gantt-chart-proto/backend/repository"
	"net/http"
)

// TODO: OpenApi定義から自動生成させる。
var rolesNeeded = map[string][]string{
	"DELETE /api/departments/:id":         {constants.RoleAdmin, constants.RoleManager},
	"DELETE /api/facilities/:id":          {constants.RoleAdmin, constants.RoleManager},
	"DELETE /api/facilitySharedLinks/:id": {constants.RoleAdmin, constants.RoleManager},
	"DELETE /api/ganttGroups/:id":         {constants.RoleAdmin, constants.RoleManager},
	"DELETE /api/holidays/:id":            {constants.RoleAdmin, constants.RoleManager},
	"DELETE /api/milestones/:id":          {constants.RoleAdmin, constants.RoleManager},
	"DELETE /api/operationSettings/:id":   {constants.RoleAdmin, constants.RoleManager},
	"DELETE /api/processes/:id":           {constants.RoleAdmin, constants.RoleManager},
	"DELETE /api/ticketUsers/:id":         {constants.RoleAdmin, constants.RoleManager},
	"DELETE /api/tickets/:id":             {constants.RoleAdmin, constants.RoleManager},
	"DELETE /api/units/:id":               {constants.RoleAdmin, constants.RoleManager},
	"DELETE /api/users/:id":               {constants.RoleAdmin, constants.RoleManager},
	"GET /api/all-tickets":                {constants.RoleAdmin, constants.RoleManager, constants.RoleWorker, constants.RoleViewer},
	"GET /api/defaultPileUps":             {constants.RoleAdmin, constants.RoleManager, constants.RoleWorker, constants.RoleViewer},
	"GET /api/departments":                {constants.RoleAdmin, constants.RoleManager, constants.RoleWorker, constants.RoleViewer, constants.RoleGuest},
	"GET /api/departments/:id":            {constants.RoleAdmin, constants.RoleManager, constants.RoleWorker, constants.RoleViewer},
	"GET /api/facilities":                 {constants.RoleAdmin, constants.RoleManager, constants.RoleWorker, constants.RoleViewer, constants.RoleGuest},
	"GET /api/facilities/:id":             {constants.RoleAdmin, constants.RoleManager, constants.RoleWorker, constants.RoleViewer, constants.RoleGuest},
	"GET /api/facilitySharedLinks":        {constants.RoleAdmin, constants.RoleManager},
	"GET /api/facilitySharedLinks/:id":    {constants.RoleAdmin, constants.RoleManager},
	"GET /api/ganttGroups":                {constants.RoleAdmin, constants.RoleManager, constants.RoleWorker, constants.RoleViewer, constants.RoleGuest},
	"GET /api/ganttGroups/:id":            {constants.RoleAdmin, constants.RoleManager, constants.RoleWorker, constants.RoleViewer, constants.RoleGuest},
	"GET /api/holidays":                   {constants.RoleAdmin, constants.RoleManager, constants.RoleWorker, constants.RoleViewer, constants.RoleGuest},
	"GET /api/holidays/:id":               {constants.RoleAdmin, constants.RoleManager, constants.RoleWorker, constants.RoleViewer},
	"GET /api/milestones":                 {constants.RoleAdmin, constants.RoleManager, constants.RoleWorker, constants.RoleViewer, constants.RoleGuest},
	"GET /api/milestones/:id":             {constants.RoleAdmin, constants.RoleManager, constants.RoleWorker, constants.RoleViewer, constants.RoleGuest},
	"GET /api/operationSettings/:id":      {constants.RoleAdmin, constants.RoleManager, constants.RoleWorker, constants.RoleViewer, constants.RoleGuest},
	"GET /api/pileUps":                    {constants.RoleAdmin, constants.RoleManager, constants.RoleWorker, constants.RoleViewer, constants.RoleGuest},
	"GET /api/processes":                  {constants.RoleAdmin, constants.RoleManager, constants.RoleWorker, constants.RoleViewer, constants.RoleGuest},
	"GET /api/processes/:id":              {constants.RoleAdmin, constants.RoleManager, constants.RoleWorker, constants.RoleViewer},
	"GET /api/scheduleAlerts":             {constants.RoleAdmin, constants.RoleManager, constants.RoleWorker, constants.RoleViewer},
	"GET /api/ticket-memo/:id":            {constants.RoleAdmin, constants.RoleManager, constants.RoleWorker, constants.RoleViewer},
	"GET /api/ticketUsers":                {constants.RoleAdmin, constants.RoleManager, constants.RoleWorker, constants.RoleViewer, constants.RoleGuest},
	"GET /api/ticketUsers/:id":            {constants.RoleAdmin, constants.RoleManager, constants.RoleWorker, constants.RoleViewer, constants.RoleGuest},
	"GET /api/tickets":                    {constants.RoleAdmin, constants.RoleManager, constants.RoleWorker, constants.RoleViewer, constants.RoleGuest},
	"GET /api/tickets/:id":                {constants.RoleAdmin, constants.RoleManager, constants.RoleWorker, constants.RoleViewer},
	"GET /api/units":                      {constants.RoleAdmin, constants.RoleManager, constants.RoleWorker, constants.RoleViewer, constants.RoleGuest},
	"GET /api/units/:id":                  {constants.RoleAdmin, constants.RoleManager, constants.RoleWorker, constants.RoleViewer},
	//"GET /api/userInfo":                   {constants.RoleAdmin, constants.RoleManager, constants.RoleWorker, constants.RoleViewer, constants.RoleGuest},
	"GET /api/users":                    {constants.RoleAdmin, constants.RoleManager, constants.RoleWorker, constants.RoleViewer, constants.RoleGuest},
	"GET /api/users/:id":                {constants.RoleAdmin, constants.RoleManager, constants.RoleWorker, constants.RoleViewer, constants.RoleGuest},
	"POST /api/copyFacilitys":           {constants.RoleAdmin, constants.RoleManager},
	"POST /api/departments":             {constants.RoleAdmin, constants.RoleManager},
	"POST /api/departments/:id":         {constants.RoleAdmin, constants.RoleManager},
	"POST /api/facilities":              {constants.RoleAdmin, constants.RoleManager},
	"POST /api/facilities/:id":          {constants.RoleAdmin, constants.RoleManager},
	"POST /api/facilitySharedLinks":     {constants.RoleAdmin, constants.RoleManager},
	"POST /api/facilitySharedLinks/:id": {constants.RoleAdmin, constants.RoleManager},
	"POST /api/ganttGroups":             {constants.RoleAdmin, constants.RoleManager, constants.RoleWorker},
	"POST /api/ganttGroups/:id":         {constants.RoleAdmin, constants.RoleManager, constants.RoleWorker},
	"POST /api/holidays":                {constants.RoleAdmin, constants.RoleManager},
	"POST /api/holidays/:id":            {constants.RoleAdmin, constants.RoleManager},
	//"POST /api/login":                     {constants.RoleAdmin, constants.RoleManager},
	"POST /api/milestones":            {constants.RoleAdmin, constants.RoleManager},
	"POST /api/milestones/:id":        {constants.RoleAdmin, constants.RoleManager},
	"POST /api/operationSettings/:id": {constants.RoleAdmin, constants.RoleManager},
	"POST /api/processes":             {constants.RoleAdmin, constants.RoleManager},
	"POST /api/processes/:id":         {constants.RoleAdmin, constants.RoleManager},
	"POST /api/ticket-memo/:id":       {constants.RoleAdmin, constants.RoleManager, constants.RoleWorker},
	"POST /api/ticketUsers":           {constants.RoleAdmin, constants.RoleManager, constants.RoleWorker},
	"POST /api/tickets":               {constants.RoleAdmin, constants.RoleManager, constants.RoleWorker},
	"POST /api/tickets/:id":           {constants.RoleAdmin, constants.RoleManager, constants.RoleWorker},
	"POST /api/units":                 {constants.RoleAdmin, constants.RoleManager},
	"POST /api/units/:id":             {constants.RoleAdmin, constants.RoleManager},
	"POST /api/users":                 {constants.RoleAdmin, constants.RoleManager},
	"POST /api/users/:id":             {constants.RoleAdmin, constants.RoleManager},
}

func getRolesFromToken(sessionID string) []string {
	userId := GetUserId(sessionID)
	userRep := repository.NewUserRepository()
	if userId != nil {
		if *userId == constants.GuestID {
			userRep = repository.NewUserRepository(repository.GuestMode)
		}
		user := userRep.Find(*userId)
		return []string{user.Role}
	}
	return []string{}
}

func RoleBasedAccessControl() gin.HandlerFunc {
	return func(c *gin.Context) {

		path := c.FullPath()
		method := c.Request.Method

		requiredRoles, ok := rolesNeeded[method+" "+path]

		if !ok {
			//c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "no permission to this resource"})
			// TODO: パスが見つからないときは認証をさせない。作りとしていまいちな気がするけど一旦速度優先。
			c.Next()
			return
		}

		sessionID, _ := c.Cookie("session_id")
		userRoles := getRolesFromToken(sessionID)

		for _, userRole := range userRoles {
			for _, requiredRole := range requiredRoles {
				if userRole == requiredRole {
					c.Next()
					return
				}
			}
		}

		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "no permission to this resource"})
	}
}
