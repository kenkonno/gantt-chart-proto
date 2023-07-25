package sub

import (
	"fmt"
	"os"
	"strings"
)

type Vuejs struct {
	StructInfo []StructInfo
}

func (v *Vuejs) GetView(structName string) string {

	file, err := os.ReadFile("vue_templates/View.vue")
	if err != nil {
		panic(err)
	}

	return RewriteString(string(file), structName)
}

func (v *Vuejs) GetComposable(structName string) string {
	file, err := os.ReadFile("vue_templates/composable.ts")
	if err != nil {
		panic(err)
	}
	result := RewriteString(string(file), structName)
	result = strings.Replace(result, "@DefaultMapping@", v.GetDefaultMapping(), -1) + "\n"
	result = strings.Replace(result, "@ResponseMapping@", v.GetResponseMapping(ToLowerCamel(structName)), -1) + "\n"
	return result
}

func (v *Vuejs) GetAsyncEdit(structName string) string {
	file, err := os.ReadFile("vue_templates/AsyncEdit.vue")
	if err != nil {
		panic(err)
	}
	result := RewriteString(string(file), structName)
	result = strings.Replace(result, "@AsyncEditMapping@", v.GetAsyncEditMapping(ToLowerCamel(structName)), -1) + "\n"
	return result
}

func (v *Vuejs) GetAsyncTable(structName string) string {
	file, err := os.ReadFile("vue_templates/AsyncTable.vue")
	if err != nil {
		panic(err)
	}
	result := RewriteString(string(file), structName)
	result = strings.Replace(result, "@AsyncTableHeader@", v.GetAsyncTableHeader(), -1) + "\n"
	result = strings.Replace(result, "@AsyncTableBody@", v.GetAsyncTableBody(), -1) + "\n"
	return result
}
func (v *Vuejs) GetAsyncTableHeader() string {
	var result []string
	for _, v := range v.StructInfo {
		result = append(result, fmt.Sprintf("        <th>%s</th>", v.Property))
	}
	return strings.Join(result, "\n")
}
func (v *Vuejs) GetAsyncTableBody() string {
	var result []string
	for _, v := range v.StructInfo {
		template := "        <td>{{ item.%s }}</td>"
		if v.Property == "Id" {
			template = `        <td @click="$emit('openEditModal', item.id)">{{ item.%s }}</td>`
		}
		result = append(result, fmt.Sprintf(template, ToSnakeCase(v.Property)))
	}
	return strings.Join(result, "\n")
}

func (v *Vuejs) GetAsyncEditMapping(lowerName string) string {
	var result []string
	var template = `    <div class="mb-2">
      <label class="form-label" for="id">%s</label>
      <input class="form-control" type="text" name="%s" id="%s" v-model="%s.%s" :disabled="%s">
    </div>
`
	for _, v := range v.StructInfo {
		pr := v.Property
		lowerPr := ToLowerCamel(pr)
		disabled := "false"
		if v.Property == "Id" || v.Property == "CreatedAt" || v.Property == "UpdatedAt" {
			disabled = "true"
		}
		result = append(result, fmt.Sprintf(template, pr, lowerPr, lowerPr, lowerName, ToSnakeCase(v.Property), disabled))
	}

	return strings.Join(result, "\n")

}

func (v *Vuejs) GetDefaultMapping() string {
	var result []string
	for _, v := range v.StructInfo {
		prop := v.Property
		value := ""

		if strings.Contains(v.Type, "int") {
			value = "0"
		}
		if strings.Contains(v.Type, "string") {
			value = `""`
		}

		if v.Property == "Id" {
			value = "null"
		}
		if v.Property == "CreatedAt" || v.Property == "UpdatedAt" {
			value = "undefined"
		}

		result = append(result, fmt.Sprintf("        %s: %s", ToSnakeCase(prop), value))
	}

	return strings.Join(result, ",\n")

}

func (v *Vuejs) GetResponseMapping(lowerName string) string {
	var result []string
	for _, v := range v.StructInfo {
		prop := v.Property
		result = append(result, fmt.Sprintf("            %s.value.%s = data.%s.%s", lowerName, ToSnakeCase(prop), lowerName, ToSnakeCase(prop)))
	}

	return strings.Join(result, "\n")

}
