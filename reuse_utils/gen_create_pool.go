package reuseutils

import (
	"strings"
)

// GenerateCreatePool 生成一个结构体使用示例代码
func GenerateCreatePool(structName string) string {
	var builder strings.Builder

	firstChar := strings.ToLower(string(structName[0]))

	builder.WriteString("\treuse.CreatePool[" + structName + "](\n")
	builder.WriteString("\t\tfunc() *" + structName + " { return &" + structName + "{} },\n")
	builder.WriteString("\t\tfunc( " + firstChar + " *" + structName + ") { " + firstChar + ".reset() },\n")
	// builder.WriteString("\t\tfunc(c *" + structName + ") { *c = " + structName + "{} },\n")
	builder.WriteString("\t\treuse.DefaultOptions[" + structName + "]().\n")
	builder.WriteString("\t\t\tWithEnableChecker())\n\n")

	return builder.String()
}

// func main() {
// 	srcFolderPath := "/Users/jensen/Projects/matrixorigin/matrixone/pkg/sql/parsers/tree"  // 设置为包含Go源代码文件的文件夹路径
// 	dstFolderPath := "/Users/jensen/demo2/generate/" // 设置为保存新生成的Go代码文件的文件夹路径

// 	err := filepath.Walk(srcFolderPath, func(path string, info os.FileInfo, err error) error {
// 		if err != nil {
// 			return err
// 		}

// 		if info.IsDir() {
// 			return nil
// 		}

// 		if filepath.Ext(path) != ".go" {
// 			return nil
// 		}

// 		// 分析Go文件
// 		fset := token.NewFileSet()
// 		node, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
// 		if err != nil {
// 			return err
// 		}

// 		var structNames []string
// 		ast.Inspect(node, func(n ast.Node) bool {
// 			// 查找结构体定义
// 			typeSpec, ok := n.(*ast.TypeSpec)
// 			if !ok {
// 				return true
// 			}

// 			_, ok = typeSpec.Type.(*ast.StructType)
// 			if ok {
// 				structNames = append(structNames, typeSpec.Name.Name)
// 			}

// 			return true
// 		})

// 		// 创建或写入目标文件
// 		dstFilePath := filepath.Join(dstFolderPath, filepath.Base(path))
// 		dstFile, err := os.Create(dstFilePath)
// 		if err != nil {
// 			return err
// 		}
// 		defer dstFile.Close()

// 		// 写入生成的代码
// 		for _, structName := range structNames {
// 			code := GenerateCreatePool(structName)
// 			_, err = dstFile.WriteString(code)
// 			if err != nil {
// 				return err
// 			}
// 		}

// 		return nil
// 	})
// 	if err != nil {
// 		panic(err)
// 	}
// }
