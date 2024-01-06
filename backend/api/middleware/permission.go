package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/kenkonno/gantt-chart-proto/backend/api/constants"
	"net/http"
)

// TODO: OpenApi定義から自動生成させる。
var rolesNeeded = map[string][]string{
	"DELETE /api/departments/:id":       {constants.RoleAdmin, constants.RoleManager},
	"DELETE /api/facilities/:id":        {constants.RoleAdmin, constants.RoleManager},
	"DELETE /api/ganttGroups/:id":       {constants.RoleAdmin, constants.RoleManager},
	"DELETE /api/holidays/:id":          {constants.RoleAdmin, constants.RoleManager},
	"DELETE /api/operationSettings/:id": {constants.RoleAdmin, constants.RoleManager},
	"DELETE /api/processes/:id":         {constants.RoleAdmin, constants.RoleManager},
	"DELETE /api/ticketUsers/:id":       {constants.RoleAdmin, constants.RoleManager},
	"DELETE /api/tickets/:id":           {constants.RoleAdmin, constants.RoleManager},
	"DELETE /api/units/:id":             {constants.RoleAdmin, constants.RoleManager},
	"DELETE /api/users/:id":             {constants.RoleAdmin, constants.RoleManager},
	"GET /api/all-tickets":              {constants.RoleAdmin, constants.RoleManager, constants.RoleWorker, constants.RoleViewer},
	"GET /api/departments":              {constants.RoleAdmin, constants.RoleManager, constants.RoleWorker, constants.RoleViewer},
	"GET /api/departments/:id":          {constants.RoleAdmin, constants.RoleManager, constants.RoleWorker, constants.RoleViewer},
	"GET /api/facilities":               {constants.RoleAdmin, constants.RoleManager, constants.RoleWorker, constants.RoleViewer},
	"GET /api/facilities/:id":           {constants.RoleAdmin, constants.RoleManager, constants.RoleWorker, constants.RoleViewer},
	"GET /api/ganttGroups":              {constants.RoleAdmin, constants.RoleManager, constants.RoleWorker, constants.RoleViewer},
	"GET /api/ganttGroups/:id":          {constants.RoleAdmin, constants.RoleManager, constants.RoleWorker, constants.RoleViewer},
	"GET /api/holidays":                 {constants.RoleAdmin, constants.RoleManager, constants.RoleWorker, constants.RoleViewer},
	"GET /api/holidays/:id":             {constants.RoleAdmin, constants.RoleManager, constants.RoleWorker, constants.RoleViewer},
	"GET /api/operationSettings/:id":    {constants.RoleAdmin, constants.RoleManager, constants.RoleWorker, constants.RoleViewer},
	"GET /api/pileUps":                  {constants.RoleAdmin, constants.RoleManager, constants.RoleWorker, constants.RoleViewer},
	"GET /api/processes":                {constants.RoleAdmin, constants.RoleManager, constants.RoleWorker, constants.RoleViewer},
	"GET /api/processes/:id":            {constants.RoleAdmin, constants.RoleManager, constants.RoleWorker, constants.RoleViewer},
	"GET /api/scheduleAlerts":           {constants.RoleAdmin, constants.RoleManager, constants.RoleWorker, constants.RoleViewer},
	"GET /api/ticketUsers":              {constants.RoleAdmin, constants.RoleManager, constants.RoleWorker, constants.RoleViewer},
	"GET /api/ticketUsers/:id":          {constants.RoleAdmin, constants.RoleManager, constants.RoleWorker, constants.RoleViewer},
	"GET /api/tickets":                  {constants.RoleAdmin, constants.RoleManager, constants.RoleWorker, constants.RoleViewer},
	"GET /api/tickets/:id":              {constants.RoleAdmin, constants.RoleManager, constants.RoleWorker, constants.RoleViewer},
	"GET /api/units":                    {constants.RoleAdmin, constants.RoleManager, constants.RoleWorker, constants.RoleViewer},
	"GET /api/units/:id":                {constants.RoleAdmin, constants.RoleManager, constants.RoleWorker, constants.RoleViewer},
	"GET /api/users":                    {constants.RoleAdmin, constants.RoleManager, constants.RoleWorker, constants.RoleViewer},
	"GET /api/users/:id":                {constants.RoleAdmin, constants.RoleManager, constants.RoleWorker, constants.RoleViewer},
	"POST /api/copyFacilitys":           {constants.RoleAdmin, constants.RoleManager},
	"POST /api/departments":             {constants.RoleAdmin, constants.RoleManager},
	"POST /api/departments/:id":         {constants.RoleAdmin, constants.RoleManager},
	"POST /api/facilities":              {constants.RoleAdmin, constants.RoleManager},
	"POST /api/facilities/:id":          {constants.RoleAdmin, constants.RoleManager},
	"POST /api/ganttGroups":             {constants.RoleAdmin, constants.RoleManager},
	"POST /api/ganttGroups/:id":         {constants.RoleAdmin, constants.RoleManager},
	"POST /api/holidays":                {constants.RoleAdmin, constants.RoleManager},
	"POST /api/holidays/:id":            {constants.RoleAdmin, constants.RoleManager},
	"POST /api/operationSettings/:id":   {constants.RoleAdmin, constants.RoleManager},
	"POST /api/processes":               {constants.RoleAdmin, constants.RoleManager},
	"POST /api/processes/:id":           {constants.RoleAdmin, constants.RoleManager},
	"POST /api/ticketUsers":             {constants.RoleAdmin, constants.RoleManager},
	"POST /api/tickets":                 {constants.RoleAdmin, constants.RoleManager},
	"POST /api/tickets/:id":             {constants.RoleAdmin, constants.RoleManager, constants.RoleWorker},
	"POST /api/units":                   {constants.RoleAdmin, constants.RoleManager},
	"POST /api/units/:id":               {constants.RoleAdmin, constants.RoleManager},
	"POST /api/users":                   {constants.RoleAdmin, constants.RoleManager},
	"POST /api/users/:id":               {constants.RoleAdmin, constants.RoleManager},
}

func getRolesFromToken(token string) []string {
	// TODO: ログイン情報からRoleを取得する
	return []string{constants.RoleAdmin}
}
func RoleBasedAccessControl() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		// TODO: ログイン情報を取得する
		userRoles := getRolesFromToken(token)

		path := c.FullPath()
		method := c.Request.Method

		requiredRoles, ok := rolesNeeded[method+" "+path]

		if !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "no permission to this resource"})
			return
		}

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
