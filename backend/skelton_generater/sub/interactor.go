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

func (r *Interactor) GetMapping(prefix string) string {

	var result string
	for _, v := range r.StructInfo {
		p := v.Property
		if p == "ID" {
			p = "Id"
		}
		value := prefix + "." + v.Property
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
			r.GetMapping("item") + `
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
` + r.GetMapping(ToLowerCamel(structName)) + `
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
		panic("invalid json")
	}
	@Lower@Rep.Upsert(db.@Upper@{
` + r.GetMapping(ToLowerCamel(structName)+"Req."+structName) + `
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
		panic("invalid json")
	}

	@Lower@Rep.Upsert(db.@Upper@{
` + r.GetMapping(ToLowerCamel(structName)+"Req."+structName) + `
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
