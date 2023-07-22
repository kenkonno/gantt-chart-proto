package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"unicode"
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
	var packageName string
	var packageMap = make(map[string]string)
	rewrite := false
	for _, v := range s {
		// impotの追加
		if strings.Contains(v, "net/http") {
			result += "\n"
			// TODO: ディレクトリ構造を変えたのでインポートを変える
			result += `@imports@
	"github.com/kenkonno/gantt-chart-proto/backend/api/openapi_models"`
			result += "\n"
		}

		// 関数内部の書き換え
		if rewrite {
			result += fmt.Sprintf(
				`	var r openapi_models.%sResponse
	r = %s.%sInvoke(c)
	c.JSON(http.StatusOK, r)
`, funcName, packageName, funcName)
			rewrite = false
		} else {
			result += v + "\n"
		}
		if strings.Contains(v, "func") {
			assigned := regexp.MustCompile(`func ([a-zA-Z]+)\(`)
			group := assigned.FindSubmatch([]byte(v))
			funcName = string(group[1])

			reg := regexp.MustCompile(`(Get|Put|Delete|Post|Invoke)`)
			packageName = ToSnakeCase(reg.ReplaceAllString(funcName, ""))
			packageMap[packageName] = packageName
			rewrite = true
		}
	}
	var imp string
	for _, v := range packageMap {
		imp += `	"github.com/kenkonno/gantt-chart-proto/backend/api/interactor/` + v + "\"\n"
	}
	result = strings.Replace(result, "@imports@", imp, -1)

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

func ToSnakeCase(s string) string {
	b := &strings.Builder{}
	for i, r := range s {
		if i == 0 {
			b.WriteRune(unicode.ToLower(r))
			continue
		}
		if unicode.IsUpper(r) {
			b.WriteRune('_')
			b.WriteRune(unicode.ToLower(r))
			continue
		}
		b.WriteRune(r)
	}
	result := b.String()
	if result == "i_d" {
		return "id"
	} else {
		return b.String()
	}
}
