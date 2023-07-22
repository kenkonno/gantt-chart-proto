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
                  $ref: '#/components/schemas/@Upper@'
              examples: {}
      operationId: get-@Lower@s
      description: ''
      parameters: []
    parameters: []
`
	result += RewriteString(template, structName)
	result += r.componentProperties()

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
			`  parameters: {}
`
	return RewriteString(template, structName)
}

func (r *Yaml) componentProperties() string {
	result := "      properties:\n"
	for _, v := range r.StructInfo {
		result += fmt.Sprintf("        %s:\n", ToSnakeCase(v.Property))
		// types
		if isNumber(v) {
			result += "          type: integer:\n"
		}
		if isString(v) {
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
