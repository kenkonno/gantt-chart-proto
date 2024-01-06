package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/samber/lo"
	"net/http"
	"os"
	"time"
)

var redisClient *redis.Client

func init() {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("SESSION_ADDR"),
		Password: os.Getenv("SESSION_PASS"), // no password set
		DB:       0,                         // use default DB
	})

	_, err := redisClient.Ping().Result()
	if err != nil {
		panic(err)
	}
}

var publicRoute = []string{
	"GET /api/userInfo",
	"POST /api/login",
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: 本来は router.Group で設定するが今回はログイン関係はここで除外する
		path := c.FullPath()
		method := c.Request.Method

		if lo.Contains(publicRoute, method+" "+path) {
			// public routeの場合は除外
			return
		}

		sessionID, _ := c.Cookie("session_id")
		if sessionID == "" { // トークンが存在しない
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
			return
		}

		if !isValidToken(sessionID) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
		}
	}
}

func isValidToken(sessionID string) bool {
	result, err := redisClient.Get(sessionID).Result()
	if err != nil {
		panic(err)
	}
	return result != ""
}

func GetUserId(sessionID string) *string {
	result, err := redisClient.Get(sessionID).Result()
	if err != nil {
		return nil
	}
	return &result
}

func UpdateSessionID(sessionID string, userId int32) {
	err := redisClient.Set(sessionID, userId, time.Hour).Err()
	if err != nil {
		panic(err)
	}
}
