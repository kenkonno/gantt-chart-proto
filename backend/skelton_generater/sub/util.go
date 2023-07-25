package sub

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"unicode"
)

func CreateFile(path string, body string) {
	create, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	err = os.WriteFile(path, []byte(body), 0777)
	if err != nil {
		panic(err)
	}
	create.Close()
}

func MakeDir(dir string) {
	unixPerms := os.ModePerm
	if err := os.MkdirAll(dir, unixPerms); err != nil {
		log.Fatal(err)
	}
}

func RewriteString(template string, value string) string {
	upper := value
	lower := ToLowerCamel(value)
	return strings.Replace(strings.Replace(template, "@Upper@", upper, -1), "@Lower@", lower, -1)
}

func ToLowerCamel(value string) string {
	return strings.ToLower(value[:1]) + value[1:]
}

func GetStructName(body []string) string {
	var r string
	for _, v := range body {
		if strings.Contains(v, "struct") {
			r = strings.Split(v, " ")[1]
		}
	}
	return r
}

func GetStructInfo(body []string) []StructInfo {
	var r []StructInfo
	start := false
	for _, v := range body {
		if strings.Contains(v, "}") {
			start = false
		}

		if start {
			re1 := regexp.MustCompile(" +")
			row := re1.ReplaceAllString(strings.Replace(v, "\t", "", -1), " ") // tabとスペース重複の削除
			arr := strings.Split(row, " ")
			if len(arr) >= 2 {
				gorm := ""
				if len(arr) == 3 {
					gorm = arr[2]
				}
				r = append(r, StructInfo{
					Property: arr[0],
					Type:     arr[1],
					Gorm:     gorm,
				})
			}
		}

		if strings.Contains(v, "struct") {
			start = true
		}
	}
	return r
}

func GetFileBody(args []string) []string {
	filepath := args[0] + ".go"
	fmt.Println(filepath)
	file, err := os.ReadFile(filepath)
	if err != nil {
		panic(err)
	}
	s := strings.Split(string(file), "\n")
	return s
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
