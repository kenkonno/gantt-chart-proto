package main

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/kenkonno/gantt-chart-proto/backend/api/constants"
	"github.com/kenkonno/gantt-chart-proto/backend/api/openapi"
	"github.com/kenkonno/gantt-chart-proto/backend/repository"
	"os"
)

var userAccountRep = repository.NewUserAccountRepository()

func main() {
	r := gin.Default()
	cfg := cors.DefaultConfig()
	cfg.AllowOrigins = []string{"http://localhost:8080"}
	cfg.AllowMethods = []string{
		"POST",
		"GET",
		"OPTIONS",
		"PUT",
		"DELETE",
	}
	cfg.AllowHeaders = []string{"Content-Type"}
	cfg.AllowCredentials = true
	cfg.AllowWildcard = true
	r.Use(cors.New(cfg))
	r.Use(userInfoMiddleware())

	r = openapi.NewRouter(r)
	r.LoadHTMLGlob("templates/*")
	// Listen and Server in 0.0.0.0:8080
	r.Run(":8081")
}

func userInfoMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie(os.Getenv(constants.TWIP_LOGIN_TOKEN_COOKIE_NAME))
		if err != nil {
			fmt.Println(err)
		}
		r := userAccountRep.FindByToken(token)
		fmt.Println(token)
		fmt.Println(r)
		if r != nil {
			c.Set("userAccountId", r.ID)
		}
	}
}
