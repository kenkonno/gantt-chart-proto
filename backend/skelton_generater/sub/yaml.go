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
	Gorm     string
}

// GetBasePaths Get/Postがデフォルト
func (r *Yaml) GetBasePaths(structName string) string {
	var result string
	template :=
		`
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
                $ref: '#/components/schemas/Get@Upper@sResponse'
              examples: {}
      operationId: get-@Lower@s
      description: ''
      parameters: []
    post:
      summary: Post@Upper@s
      operationId: post-@Lower@s
      responses:
        '200':
          description: OK
      description: ''
      requestBody:
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/Post@Upper@sRequest'
    parameters: []
`
	result += RewriteString(template, structName)

	return result
}

func (r *Yaml) GetWithIdPath(structName string) string {
	var result string
	template :=
		`
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
    delete:
      summary: Delete@Upper@sId
      operationId: delete-@Lower@s-id
      responses:
        '200':
          description: OK
      description: ''
    post:
      summary: Post@Upper@sId
      operationId: post-@Lower@s-id
      responses:
        '200':
          description: OK
      requestBody:
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/Post@Upper@sRequest'
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
	var requiredProps []string
	for _, v := range r.StructInfo {
		result += fmt.Sprintf("        %s:\n", ToSnakeCase(v.Property))
		// types
		if isUpdatedAt(v) {
			result += "          type: int\n"
		} else if isNumber(v) {
			result += "          type: integer\n"
		} else if isString(v) || isDatetime(v) {
			result += "          type: string\n"
		}

		// formats
		if isID(v) {
			result += "          format: int32\n"
		}
		if isDatetime(v) {
			result += "          format: date-time\n"
		}
		// Null, required関連の判定
		if !strings.Contains(v.Type, "*") {
			if v.Property == "UpdatedAt" || v.Property == "CreatedAt" || v.Property == "Date" { // TODO: この辺り整理しないといけない。Dateでmin=1だとvalidationエラーになるから一旦回避
			} else {
				requiredProps = append(requiredProps, ToSnakeCase(v.Property))
				// カスタムはこんな感じで追加する。
				// x-go-custom-tag: binding:"min=1"
				// requiredの時は 1以上として簡易的にバリデーションとする。
				result += "          x-go-custom-tag: binding:\"min=1\"\n"
			}
		} else {
			result += "          nullable: true\n"
		}
	}
	if len(requiredProps) > 0 {
		result += "      required:\n"
		for _, v := range requiredProps {
			result += fmt.Sprintf("        - %s\n", v)
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
        @Lower@:
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

func (r *Yaml) GetDeleteIdRequest(structName string) string {
	template :=
		`    Delete@Upper@sIdRequest:
      title: Delete@Upper@sIdRequest
      type: object
      properties:
        id:
          type: integer
          format: int32
`
	return RewriteString(template, structName)
}

func (r *Yaml) GetDeleteIdResponse(structName string) string {
	template := `    Delete@Upper@sIdResponse:
      title: Delete@Upper@sIdResponse
      type: object
      properties:
        msg:
          type: string
`
	return RewriteString(template, structName)
}

func (r *Yaml) GetPostRequest(structName string) string {
	template := `    Post@Upper@sRequest:
      title: Post@Upper@sRequest
      type: object
      properties:
        @Lower@:
          $ref: '#/components/schemas/@Upper@'
`
	return RewriteString(template, structName)
}

func (r *Yaml) GetPostResponse(structName string) string {
	template := `    Post@Upper@sResponse:
      title: Post@Upper@sResponse
      type: object
      properties:
        @Lower@:
          $ref: '#/components/schemas/@Upper@'
`
	return RewriteString(template, structName)
}

func (r *Yaml) GetPostIdRequest(structName string) string {
	template := `    Post@Upper@sIdRequest:
      title: Post@Upper@sIdRequest
      type: object
      properties:
        @Lower@:
          $ref: '#/components/schemas/@Upper@'
`
	return RewriteString(template, structName)
}

func (r *Yaml) GetPostIdResponse(structName string) string {
	template := `    Post@Upper@sIdResponse:
      title: Post@Upper@sIdResponse
      type: object
      properties:
        @Lower@:
          $ref: '#/components/schemas/@Upper@'
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

func isUpdatedAt(s StructInfo) bool {
	if strings.Contains(s.Property, "UpdatedAt") {
		return true
	}
	return false
}
