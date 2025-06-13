package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/kenkonno/gantt-chart-proto/backend/api/middleware"
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
	r.Use(middleware.RoleBasedAccessControl())
	r.Use(middleware.AuthMiddleware())
	r.Use(middleware.GuestAccessControl())
	// TODO: FeatureOptionの認可を追加する。この仕組みだとAPIのパスを見てDBのオプション有無で認可するって感じかな。そのうち量が多くなると必要になる。
	r.Use(gzip.Gzip(gzip.DefaultCompression))

	r = openapi.NewRouter(r)
	//r.LoadHTMLGlob("templates/*") TODO: たぶん現状では不要。
	// Listen and Server in 0.0.0.0:8080
	r.Run(":" + os.Getenv("API_PORT"))
}
