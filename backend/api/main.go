package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/kenkonno/gantt-chart-proto/backend/api/openapi"
	"os"
)

func main() {
	r := gin.Default()
	cfg := cors.DefaultConfig()
	//cfg.AllowOrigins = []string{"http://localhost:8080"}
	cfg.AllowOrigins = []string{os.Getenv("ORIGIN")}
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
	//r.Use(userInfoMiddleware())

	r = openapi.NewRouter(r)
	r.LoadHTMLGlob("templates/*")
	// Listen and Server in 0.0.0.0:8080
	r.Run(":" + os.Getenv("API_PORT"))
}

//func userInfoMiddleware() gin.HandlerFunc {
//	return func(c *gin.Context) {
//		token, err := c.Cookie(os.Getenv(constants.TWIP_LOGIN_TOKEN_COOKIE_NAME))
//		if err != nil {
//			fmt.Println(err)
//		}
//		r := userAccountRep.FindByToken(token)
//		fmt.Println(token)
//		fmt.Println(r)
//		if r != nil {
//			c.Set("userAccountId", r.ID)
//		}
//	}
//}
