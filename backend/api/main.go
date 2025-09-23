package main

import (
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/kenkonno/gantt-chart-proto/backend/api/middleware"
	"github.com/kenkonno/gantt-chart-proto/backend/api/openapi"
)

func main() {
	r := gin.Default()

	// CORSの設定を準備
	cfg := cors.Config{
		AllowCredentials: true,
		AllowMethods: []string{
			"POST", "GET", "OPTIONS", "PUT", "DELETE",
		},
		AllowHeaders: []string{
			"Access-Control-Allow-Headers", "Content-Type", "Content-Length",
			"Accept-Encoding", "X-CSRF-Token", "Authorization",
		},
	}

	// 環境変数 ORIGIN の値に応じて設定を分岐
	origin := os.Getenv("ORIGIN")
	if origin == "*" {
		// "*" の場合は、どのオリジンも動的に許可する
		// これにより認証情報 (Credentials) との併用が可能になります
		cfg.AllowOriginFunc = func(origin string) bool {
			return true
		}
	} else {
		// "*" 以外の場合は、指定されたオリジンのみを許可する
		cfg.AllowOrigins = []string{origin}
	}

	r.Use(cors.New(cfg))
	r.Use(middleware.RoleBasedAccessControl())
	r.Use(middleware.AuthMiddleware())
	r.Use(middleware.GuestAccessControl())
	r.Use(gzip.Gzip(gzip.DefaultCompression))

	r = openapi.NewRouter(r)
	r.Run(":" + os.Getenv("API_PORT"))
}