package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"strings"

	reuseutils "main/reuse_utils"
)

// 指定生成的代码类型
const (
	srcFolderPath = "/Users/jensen/Projects/matrixorigin/matrixone/pkg/sql/parsers/tree/alter.go" // 设置为包含Go源代码文件的文件夹路径

	inValid = iota
	genCreatePool
	genTypeName
	genReset
	genFree
	format
)

// 请手动指定
var genType = inValid

var (
	dstFileName   = "generate"
	generate      func(*ast.StructType, string) string
	dstFolderPath string
)

func init() {
	switch genType {
	case genCreatePool:
		generate = reuseutils.GenerateCreatePool
	case genTypeName:
		generate = reuseutils.GenerateTypeName
	case genReset:
		generate = reuseutils.GenerateReset
	case genFree:
		generate = reuseutils.GenerateFree
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
	dstFolderPath = currentDir + "/" + dstFileName
	err = os.MkdirAll(dstFolderPath, 0o755)
	if err != nil {
		fmt.Println("无法创建文件夹：", err)
		return
	}
}

func main() {
	var err error
	if genType == format {
		err = filepath.Walk(srcFolderPath, Format)
	} else {
		err = filepath.Walk(srcFolderPath, Generate)
	}
	if err != nil {
		panic(err)
	}
}

func Generate(path string, info os.FileInfo, err error) error {
	if err != nil {
		return err
	}

	if info.IsDir() || filepath.Ext(path) != ".go" {
		return nil
	}

	// Parse the Go source file.
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
	if err != nil {
		return err
	}

	var structNames []string
	var structTypes []*ast.StructType
	ast.Inspect(node, func(n ast.Node) bool {
		if typeSpec, isTypeSpec := n.(*ast.TypeSpec); isTypeSpec {
			if structType, isStruct := typeSpec.Type.(*ast.StructType); isStruct {
				if reuseutils.HasFormatMethod(node, typeSpec.Name.Name) {
					// GenTypeName and CreatePool only need structName
					structNames = append(structNames, typeSpec.Name.Name)
					switch genType {
					case genReset:
						// 如果存在Format方法，则生成reset方法
						structTypes = append(structTypes, structType)
					}
				}
			}
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
	for i, structName := range structNames {
		var code string

		switch genType {
		case genReset:
			code = generate(structTypes[i], structName)
		default:
			code = generate(nil, structName)
		}

		_, err = dstFile.WriteString(code)
		if err != nil {
			return err
		}
	}

	return nil
}

func Format(fp string, fi os.FileInfo, err error) error {
	if err != nil || fi.IsDir() || strings.HasSuffix(fp, "_test.go") {
		return err
	}
	if strings.HasSuffix(fi.Name(), ".go") {
		log.Printf("Processing file: %s", fp)
		err := reuseutils.FormatFile(fp)
		if err != nil {
			log.Printf("Error processing file '%s': %v", fp, err)
			return err
		}
	}
	return nil
}
