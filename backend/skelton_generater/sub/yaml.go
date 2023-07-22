package sub

import (
	"fmt"
	"strings"
)

type Yaml struct {
	StructInfo []StructInfo
}

type StructInfo struct {
	Property string
	Type     string
}

// GetPaths
func (r *Yaml) GetGetPaths(structName string) string {
	var result string
	template :=
		`paths:
  /api/@Lower@s:
    get:
      summary: Get@Upper@s
      tags: []
      responses:
        '200':
          description: ''
          content:
            application/json:
              schema:
                type: array
                items:
                  $ref: '#/components/schemas/Get@Upper@sResponse'
              examples: {}
      operationId: get-@Lower@s
      description: ''
      parameters: []
    parameters: []
`
	result += RewriteString(template, structName)

	return result
}

func (r *Yaml) GetGetIdPaths(structName string) string {
	var result string
	template :=
		`paths:
  /api/@Lower@s/{id}:
    parameters:
      - schema:
          type: number
        name: id
        in: path
        required: true
    get:
      summary: Get@Upper@sId
      tags: []
      responses:
        '200':
          description: ''
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Get@Upper@sIdResponse'
              examples: {}
      operationId: get-@Lower@s-id
      description: ''
      parameters: []
`
	result += RewriteString(template, structName)

	return result
}
func (r *Yaml) GetComponents(structName string) string {
	template :=
		`components:
  schemas:
    @Upper@:
      title: @Upper@
      type: object
` + r.componentProperties() +
			`      description: ''
      x-tags:
        - @Upper@
`
	return RewriteString(template, structName)
}

func (r *Yaml) componentProperties() string {
	result := "      properties:\n"
	for _, v := range r.StructInfo {
		result += fmt.Sprintf("        %s:\n", ToSnakeCase(v.Property))
		// types
		if isNumber(v) {
			result += "          type: integer\n"
		}
		if isString(v) || isDatetime(v) {
			result += "          type: string\n"
		}

		// formats
		if isID(v) {
			result += "          format: int32\n"
		}
		if isDatetime(v) {
			result += "          format: date-time\n"
		}
	}
	return result
}

func (r *Yaml) GetGetResponse(structName string) string {
	template :=
		`    Get@Upper@sResponse:
      title: Get@Upper@sResponse
      type: object
      properties:
        list:
          type: array
          items:
            $ref: '#/components/schemas/@Upper@'
      required:
        - list
`
	return RewriteString(template, structName)
}
func (r *Yaml) GetGetIdResponse(structName string) string {
	template :=
		`    Get@Upper@sIdResponse:
      title: Get@Upper@sResponse
      type: object
      properties:
        user:
          $ref: '#/components/schemas/@Upper@'
`
	return RewriteString(template, structName)
}
func (r *Yaml) GetGetRequest(structName string) string {
	template :=
		`    Get@Upper@sRequest:
      title: Get@Upper@sRequest
      type: object
`
	return RewriteString(template, structName)
}

func (r *Yaml) GetGetIdRequest(structName string) string {
	template :=
		`    Get@Upper@sIdRequest:
      title: Get@Upper@sIdRequest
      type: object
      properties:
        id:
          type: integer
          format: int32
`
	return RewriteString(template, structName)
}

func isNumber(s StructInfo) bool {
	if strings.Contains(s.Type, "int") {
		return true
	}
	return false
}

func isID(s StructInfo) bool {
	if strings.Contains(s.Property, "ID") {
		return true
	}
	return false
}

func isString(s StructInfo) bool {
	if strings.Contains(s.Type, "string") {
		return true
	}
	return false
}

func isDatetime(s StructInfo) bool {
	if strings.Contains(s.Type, "time.Time") {
		return true
	}
	return false
}
