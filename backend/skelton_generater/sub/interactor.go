package sub

import (
	"fmt"
	"strings"
)

type Interactor struct {
	StructInfo []StructInfo
}

func (r *Interactor) GetPackage(structName string) string {
	return "package " + ToSnakeCase(structName) + "s\n"
}

func (r *Interactor) GetImports() string {
	return `import (
	"github.com/gin-gonic/gin"
	"github.com/kenkonno/gantt-chart-proto/backend/api/openapi_models"
	"github.com/kenkonno/gantt-chart-proto/backend/repository"
)
`
}

func (r *Interactor) GetMapping(prefix string, withoutId bool) string {

	var result string
	for _, v := range r.StructInfo {

		// プロパティ定義（左側）
		p := v.Property
		if p == "Id" && withoutId {
			continue
		}
		// value定義（右側）
		value := prefix + "." + v.Property

		if v.Property == "UpdatedAt" {
			value = "0"
		}
		if v.Property == "CreatedAt" {
			value = "time.Time{}"
		}

		result += fmt.Sprintf("				%s:        %s,\n", p, value)
	}

	return result

}

// GetIdInvoke Get With Id
func (r *Interactor) GetInvoke(structName string) string {
	template :=
		`func @Method@@Upper@sInvoke(c *gin.Context) openapi_models.@Method@@Upper@sResponse {
	@Lower@Rep := repository.New@Upper@Repository()

	@Lower@List := @Lower@Rep.FindAll()

	return openapi_models.@Method@@Upper@sResponse{
		List: lo.Map(@Lower@List, func(item db.@Upper@, index int) openapi_models.@Upper@ {
			return openapi_models.@Upper@{
` +
			r.GetMapping("item", false) + `
			}
		}),
	}
}
`
	return strings.Replace(RewriteString(template, structName), "@Method@", "Get", -1)
}

// GetIdInvoke Get With Id
func (r *Interactor) GetIdInvoke(structName string) string {
	template :=
		`func @Method@@Upper@sIdInvoke(c *gin.Context) openapi_models.@Method@@Upper@sIdResponse {
	@Lower@Rep := repository.New@Upper@Repository()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		panic(err)
	}

	@Lower@ := @Lower@Rep.Find(int32(id))

	return openapi_models.GetUsersIdResponse{
		User: openapi_models.User{
` + r.GetMapping(ToLowerCamel(structName), false) + `
		},
	}
}
`
	return strings.Replace(RewriteString(template, structName), "@Method@", "Get", -1)

}

func (r *Interactor) PostInvoke(structName string) string {
	template :=
		`func @Method@@Upper@sInvoke(c *gin.Context) openapi_models.@Method@@Upper@sResponse {

	@Lower@Rep := repository.New@Upper@Repository()

	var @Lower@Req openapi_models.@Method@@Upper@sRequest
	if err := c.ShouldBindJSON(&@Lower@Req); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		panic(err)
	}
	@Lower@Rep.Upsert(db.@Upper@{
` + r.GetMapping(ToLowerCamel(structName)+"Req."+structName, true) + `
	})

	return openapi_models.@Method@@Upper@sResponse{}

}
`
	return strings.Replace(RewriteString(template, structName), "@Method@", "Post", -1)

}

// GetIdInvoke Get With Id
func (r *Interactor) PostIdInvoke(structName string) string {
	template :=
		`func @Method@@Upper@sIdInvoke(c *gin.Context) openapi_models.@Method@@Upper@sIdResponse {

	@Lower@Rep := repository.New@Upper@Repository()

	var @Lower@Req openapi_models.@Method@@Upper@sRequest
	if err := c.ShouldBindJSON(&@Lower@Req); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		panic(err)
	}

	@Lower@Rep.Upsert(db.@Upper@{
` + r.GetMapping(ToLowerCamel(structName)+"Req."+structName, false) + `
	})

	return openapi_models.@Method@@Upper@sIdResponse{}

}

`
	return strings.Replace(RewriteString(template, structName), "@Method@", "Post", -1)

}

// GetIdInvoke Get With Id
func (r *Interactor) DeleteIdInvoke(structName string) string {
	template :=
		`func @Method@@Upper@sIdInvoke(c *gin.Context) openapi_models.@Method@@Upper@sIdResponse {

	@Lower@Rep := repository.New@Upper@Repository()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		panic(err)
	}

	@Lower@Rep.Delete(int32(id))

	return openapi_models.@Method@@Upper@sIdResponse{}

}
`
	return strings.Replace(RewriteString(template, structName), "@Method@", "Delete", -1)

}
