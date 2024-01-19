package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"

	reuseutils "main/reuse_utils"
)

// 指定路径
var (
	srcFolderPath = "/Users/jensen/Projects/matrixorigin/matrixone/pkg/sql/parsers/tree" // 设置为包含Go源代码文件的文件夹路径
	dstFileName   = "generate"
)

// 指定生成的代码类型
const (
	InValid = iota
	GenCreatePool
	GenTypeName
	GenRelease
)

// 请手动指定
var genType = GenCreatePool

func main() {
	var generate func(string) string

	switch genType {
	case GenCreatePool:
		generate = reuseutils.GenerateCreatePool
	case GenTypeName:
		generate = reuseutils.GenerateTypeName
	// case GenCreatePool:
	// 	generate = reuseutils.GenerateCreatePool	case GenCreatePool:
	default:
		fmt.Println("Please specify the type of code to generate.")
		os.Exit(1)
	}

	// 获取当前工作目录
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Println("无法获取当前工作目录：", err)
		return
	}
	dstFolderPath := currentDir + "/" + dstFileName
	err = os.MkdirAll(dstFolderPath, 0o755)
	if err != nil {
		fmt.Println("无法创建文件夹：", err)
		return
	}

	err = filepath.Walk(srcFolderPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		if filepath.Ext(path) != ".go" {
			return nil
		}

		// 分析Go文件
		fset := token.NewFileSet()
		node, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
		if err != nil {
			return err
		}

		var structNames []string
		ast.Inspect(node, func(n ast.Node) bool {
			// 查找结构体定义
			typeSpec, ok := n.(*ast.TypeSpec)
			if !ok {
				return true
			}

			_, ok = typeSpec.Type.(*ast.StructType)
			if ok {
				structNames = append(structNames, typeSpec.Name.Name)
			}
			return true
		})

		// 创建或写入目标文件
		dstFilePath := filepath.Join(dstFolderPath, filepath.Base(path))
		dstFile, err := os.Create(dstFilePath + "_")
		if err != nil {
			return err
		}
		defer dstFile.Close()

		// 写入生成的代码
		for _, structName := range structNames {
			code := generate(structName)
			_, err = dstFile.WriteString(code)
			if err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		panic(err)
	}
}
