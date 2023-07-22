package sub

import "strings"

type Interactor struct {
}

func (r *Interactor) GetPackage(structName string) string {
	return "package " + ToSnakeCase(structName) + "\n"
}

func (r *Interactor) GetImports() string {
	return `import (
	"github.com/gin-gonic/gin"
	"github.com/kenkonno/gantt-chart-proto/backend/api/openapi_models"
	"github.com/kenkonno/gantt-chart-proto/backend/repository"
)
`
}

// GetIdInvoke Get With Id
func (r *Interactor) GetInvoke(structName string) string {
	template :=
		`func @Method@@Upper@sInvoke(c *gin.Context) openapi_models.@Method@@Upper@sResponse {
	@Lower@Rep := repository.New@Upper@Repository()

	@Lower@List := @Lower@Rep.FindAll()

	return openapi_models.@Method@@Upper@sResponse{}
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

	return openapi_models.@Method@@Upper@sIdResponse{}
}
`
	return strings.Replace(RewriteString(template, structName), "@Method@", "Get", -1)

}

// GetIdInvoke Get With Id
func (r *Interactor) PostInvoke(structName string) string {
	template :=
		`func @Method@@Upper@sInvoke(c *gin.Context) openapi_models.@Method@@Upper@sResponse {

	@Lower@Rep := repository.New@Upper@Repository()

	var @Lower@Req openapi_models.@Method@@Upper@sRequest
	if err := c.ShouldBindJSON(&@Lower@Req); err != nil {
		panic("invalid json")
	}

	@Lower@Rep.Upsert(db.@Upper@{
		UpdatedAt:     0,
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

	var @Lower@Req openapi_models.@Method@@Upper@sIdRequest
	if err := c.ShouldBindJSON(&@Lower@Req); err != nil {
		panic("invalid json")
	}

	@Lower@Rep.Upsert(db.@Upper@{
		ID: @Lower@Req.ID,
		UpdatedAt:     0,
	})

	return openapi_models.@Method@@Upper@sIdResponse{}

}
`
	return strings.Replace(RewriteString(template, structName), "@Method@", "Post", -1)

}

// GetIdInvoke Get With Id
func (r *Interactor) DeleteInvoke(structName string) string {
	template :=
		`func @Method@@Upper@sInvoke(c *gin.Context) openapi_models.@Method@@Upper@sResponse {

	@Lower@Rep := repository.New@Upper@Repository()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		panic(err)
	}

	@Lower@Rep.Delete(int32(id))

	return openapi_models.@Method@@Upper@sResponse{}

}
`
	return strings.Replace(RewriteString(template, structName), "@Method@", "Delete", -1)

}
