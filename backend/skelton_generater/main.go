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

	// post with Id
	interactorPostResult := i.GetPackage(structName)
	interactorPostResult += i.GetImports()
	interactorPostResult += i.PostInvoke(structName)
	sub.CreateFile("dest/interactor/"+sub.ToSnakeCase(structName)+"s/"+"post_"+sub.ToSnakeCase(structName+"sId")+".go", interactorPostResult)

	// delete
	interactorDeleteResult := i.GetPackage(structName)
	interactorDeleteResult += i.GetImports()
	interactorDeleteResult += i.DeleteInvoke(structName)
	sub.CreateFile("dest/interactor/"+sub.ToSnakeCase(structName)+"s/"+"delete_"+sub.ToSnakeCase(structName)+"s.go", interactorDeleteResult)

	// yaml
	y := sub.Yaml{structInfo}
	yamlResult := "########## PATHS ################### \n\n"
	yamlResult += y.GetGetPaths(structName)
	yamlResult += y.GetGetIdPaths(structName)
	yamlResult += "########## MODELS(components) ################### \n\n"
	yamlResult += y.GetComponents(structName)
	yamlResult += y.GetGetRequest(structName) // TODO: Request周りは結局Post系しか使ってないけど一応作る
	yamlResult += y.GetGetResponse(structName)
	yamlResult += y.GetGetIdRequest(structName) // TODO: Request周りは結局Post系しか使ってないけど一応作る
	yamlResult += y.GetGetIdResponse(structName)
	yamlResult += y.GetDeleteRequest(structName)
	yamlResult += y.GetDeleteResponse(structName)
	sub.CreateFile("dest/yaml_info.txt", yamlResult)

	fmt.Println("完了")

}
