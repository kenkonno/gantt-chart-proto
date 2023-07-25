package main

import (
	"flag"
	"fmt"
	"github.com/kenkonno/skelton_generator/sub"
)

// 使い方
func main() {

	flag.Parse()
	args := flag.Args()

	fmt.Println("開始")
	// ファイルの読み込み
	fileBody := sub.GetFileBody(args)
	structName := sub.GetStructName(fileBody)
	lowerName := sub.ToLowerCamel(structName)
	structInfo := sub.GetStructInfo(fileBody)
	// destディレクトリの作成
	sub.MakeDir("dest")
	sub.MakeDir("dest/repository")
	sub.MakeDir("dest/interactor/" + sub.ToSnakeCase(structName) + "s")

	r := sub.Repository{}

	// repositoryファイルの作成
	var repositoryResult string
	repositoryResult += r.GetPackage()
	repositoryResult += r.GetImports()
	repositoryResult += "// Auto generated start \n"
	repositoryResult += r.GetConstructor(structName)
	repositoryResult += r.GetDefaultFunctions(structName)
	repositoryResult += "// Auto generated end \n"
	sub.CreateFile("dest/repository/"+sub.ToSnakeCase(structName)+".go", repositoryResult)

	// interactor の作成
	// get
	i := sub.Interactor{StructInfo: structInfo}
	interactorGetResult := i.GetPackage(structName)
	interactorGetResult += i.GetImports()
	interactorGetResult += i.GetInvoke(structName)
	sub.CreateFile("dest/interactor/"+sub.ToSnakeCase(structName)+"s/"+"get_"+sub.ToSnakeCase(structName)+"s.go", interactorGetResult)

	// get with id
	interactorGetIdResult := i.GetPackage(structName)
	interactorGetIdResult += i.GetImports()
	interactorGetIdResult += i.GetIdInvoke(structName)
	sub.CreateFile("dest/interactor/"+sub.ToSnakeCase(structName)+"s/"+"get_"+sub.ToSnakeCase(structName+"sId")+".go", interactorGetIdResult)

	// post
	interactorPostResult := i.GetPackage(structName)
	interactorPostResult += i.GetImports()
	interactorPostResult += i.PostInvoke(structName)
	sub.CreateFile("dest/interactor/"+sub.ToSnakeCase(structName)+"s/"+"post_"+sub.ToSnakeCase(structName+"s")+".go", interactorPostResult)

	// post with Id TODO: 変える必要あるかも
	interactorPostIdResult := i.GetPackage(structName)
	interactorPostIdResult += i.GetImports()
	interactorPostIdResult += i.PostIdInvoke(structName)
	sub.CreateFile("dest/interactor/"+sub.ToSnakeCase(structName)+"s/"+"post_"+sub.ToSnakeCase(structName+"sId")+".go", interactorPostIdResult)

	// delete with Id
	interactorDeleteResult := i.GetPackage(structName)
	interactorDeleteResult += i.GetImports()
	interactorDeleteResult += i.DeleteIdInvoke(structName)
	sub.CreateFile("dest/interactor/"+sub.ToSnakeCase(structName)+"s/"+"delete_"+sub.ToSnakeCase(structName+"sId")+"s.go", interactorDeleteResult)

	// yaml
	// mustache:https://github.com/OpenAPITools/openapi-generator/blob/master/modules/openapi-generator/src/main/resources/go-gin-server/model.mustache
	y := sub.Yaml{structInfo}
	yamlResult := "########## PATHS ################### \n\n"
	yamlResult += y.GetBasePaths(structName)
	yamlResult += y.GetWithIdPath(structName)
	yamlResult += "########## MODELS(components) ################### \n\n"
	// Model
	yamlResult += y.GetComponents(structName)
	// Get
	yamlResult += y.GetGetRequest(structName) // TODO: Request周りは結局Post系しか使ってないけど一応作る
	yamlResult += y.GetGetResponse(structName)
	// GetId
	yamlResult += y.GetGetIdRequest(structName) // TODO: Request周りは結局Post系しか使ってないけど一応作る
	yamlResult += y.GetGetIdResponse(structName)
	// Delete
	yamlResult += y.GetDeleteIdRequest(structName)
	yamlResult += y.GetDeleteIdResponse(structName)
	// Post
	yamlResult += y.GetPostRequest(structName)
	yamlResult += y.GetPostResponse(structName)

	// Post
	yamlResult += y.GetPostIdRequest(structName)
	yamlResult += y.GetPostIdResponse(structName)
	sub.CreateFile("dest/yaml_info.txt", yamlResult)

	// Vuejs
	v := sub.Vuejs{StructInfo: structInfo}
	// View
	viewResult := v.GetView(structName)
	sub.MakeDir("dest/src/views")
	sub.CreateFile("dest/src/views/"+structName+"View.vue", viewResult)

	// Composable
	composableResult := v.GetComposable(structName)
	sub.MakeDir("dest/src/composable")
	sub.CreateFile("dest/src/composable/"+lowerName+".ts", composableResult)

	// Edit
	asyncEditResult := v.GetAsyncEdit(structName)
	sub.MakeDir("dest/src/components/" + lowerName)
	sub.CreateFile("dest/src/components/"+lowerName+"/Async"+structName+"Edit.vue", asyncEditResult)

	// Table
	asyncTableResult := v.GetAsyncTable(structName)
	sub.MakeDir("dest/src/components/" + lowerName)
	sub.CreateFile("dest/src/components/"+lowerName+"/Async"+structName+"Table.vue", asyncTableResult)

	fmt.Println("完了")

}
