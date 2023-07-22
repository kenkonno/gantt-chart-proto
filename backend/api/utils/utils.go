package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-module/carbon/v2"
	"strings"
)

var DateTimeFormat string = "Y-m-d H:i:s"
var DateTimeFormatUTC string = "Y-m-d H:i:s +00:00"

func StrToDate(dateStr string) carbon.Carbon {
	return carbon.Parse(dateStr)
}

func SplitIds(ids []string) [][]string {

	var results [][]string
	var workRow []string

	for i, v := range ids {
		workRow = append(workRow, v)
		if i%20 == 0 {
			results = append(results, workRow)
			workRow = []string{}
		}
	}
	if len(workRow) > 0 {
		results = append(results, workRow)
	}
	return results
}

func Truncate(c carbon.Carbon) carbon.Carbon {
	return c.SetHour(0).SetMinute(0).SetSecond(0).SetMillisecond(0)
}

func Now29() carbon.Carbon {
	c := carbon.Now("Asia/Tokyo")
	if 0 <= c.Hour() && c.Hour() < 5 {

		return Truncate(c.Yesterday())
	}
	return Truncate(carbon.Now("Asia/Tokyo"))
}

func ReplaceThumbnail(url string) string {
	url = strings.Replace(url, "{width}", "128", 1)
	url = strings.Replace(url, "{height}", "128", 1)
	return url
}

func GetUserAccountId(c *gin.Context) string {
	var token string
	if v, ok := c.Get("userAccountId"); ok {
		token = v.(string)
	} else {
		return ""
	}
	return token
}
