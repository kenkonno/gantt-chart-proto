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
	sub.MakeDir("dest/interactor/" + sub.ToSnakeCase(structName))

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
	i := sub.Interactor{}
	interactorGetResult := i.GetPackage(structName)
	interactorGetResult += i.GetImports()
	interactorGetResult += i.GetInvoke(structName)
	sub.CreateFile("dest/interactor/"+sub.ToSnakeCase(structName)+"/"+"get_"+sub.ToSnakeCase(structName)+"s.go", interactorGetResult)

	// get with id
	interactorGetIdResult := i.GetPackage(structName)
	interactorGetIdResult += i.GetImports()
	interactorGetIdResult += i.GetIdInvoke(structName)
	sub.CreateFile("dest/interactor/"+sub.ToSnakeCase(structName)+"/"+"get_"+sub.ToSnakeCase(structName+"Id")+"s.go", interactorGetIdResult)

	// post with Id
	interactorPostResult := i.GetPackage(structName)
	interactorPostResult += i.GetImports()
	interactorPostResult += i.PostInvoke(structName)
	sub.CreateFile("dest/interactor/"+sub.ToSnakeCase(structName)+"/"+"post_"+sub.ToSnakeCase(structName+"Id")+"s.go", interactorPostResult)

	// delete
	interactorDeleteResult := i.GetPackage(structName)
	interactorDeleteResult += i.GetImports()
	interactorDeleteResult += i.DeleteInvoke(structName)
	sub.CreateFile("dest/interactor/"+sub.ToSnakeCase(structName)+"/"+"delete_"+sub.ToSnakeCase(structName)+"s.go", interactorDeleteResult)

	// yaml
	y := sub.Yaml{structInfo}
	yamlResult := y.GetGetPaths(structName)
	yamlResult += y.GetComponents(structName)
	sub.CreateFile("dest/yaml_info.txt", yamlResult)

	fmt.Println("完了")

}
