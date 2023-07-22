package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

func main() {
	fmt.Println("開始")
	// ファイルの読み込み
	file, err := os.ReadFile("openapi/api_default.go")
	if err != nil {
		panic(err)
	}
	s := strings.Split(string(file), "\n")

	var result string
	var funcName string
	rewrite := false
	for _, v := range s {
		// impotの追加
		if strings.Contains(v, "net/http") {
			result += "\n"
			// TODO: ディレクトリ構造を変えたのでインポートを変える
			result += `	"github.com/kenkonno/gantt-chart-proto/backend/api/interactor"
	"github.com/kenkonno/gantt-chart-proto/backend/api/openapi_models"`
			result += "\n"
		}

		// 関数内部の書き換え
		if rewrite {
			result += fmt.Sprintf(
				`	var r openapi_models.%sResponse
	r = interactor.%sInvoke(c)
	c.JSON(http.StatusOK, r)
`, funcName, funcName)
			rewrite = false
		} else {
			result += v + "\n"
		}
		if strings.Contains(v, "func") {
			assigned := regexp.MustCompile(`func ([a-zA-Z]+)\(`)
			group := assigned.FindSubmatch([]byte(v))
			funcName = string(group[1])
			rewrite = true
		}
	}

	create, err := os.Create("tmp_api_default.go")
	if err != nil {
		panic(err)
	}
	err = os.WriteFile("tmp_api_default.go", []byte(result), 0777)
	if err != nil {
		panic(err)
	}

	err = os.Remove("openapi/api_default.go")
	if err != nil {
		panic(err)
	}

	create.Close()

	err = os.Rename("tmp_api_default.go", "openapi/api_default.go")
	if err != nil {
		panic(err)
	}
	fmt.Println("完了")

}
