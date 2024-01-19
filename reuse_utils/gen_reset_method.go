package reuseutils

import (
	"go/ast"
	"os"
	"strings"
)

func getElementType(expr ast.Expr) (string, bool) {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name, true // 直接的类型名称
	case *ast.StarExpr:
		return expr.(*ast.StarExpr).X.(*ast.Ident).Name, true // 数组成员也是指针
	// ...可以添加更多的case来处理更复杂的类型情况
	default:
		return "", false // 未知或不支持的类型
	}
}

// GenerateResetFunc 生成结构体的 reset 方法
func GenerateResetFunc(structType *ast.StructType, structName string) string {
	var builder strings.Builder
	builder.WriteString("\nfunc (node *" + structName + ") reset() {\n")

	for _, field := range structType.Fields.List {
		// 对于类型为 *ast.ArrayType (即 slice 类型)，我们需要遍历这个 slice 并释放内容。
		if arrayType, ok := field.Type.(*ast.ArrayType); ok {
			// 使用类型断言来进一步断言被释放的对象类型, 这里其实会比较麻烦..

			eltTypeName, ok := getElementType(arrayType.Elt)
			if !ok {
				// 无法处理这个元素类型, 打印它并且退出
				ast.Fprint(os.Stdout, nil, arrayType.Elt, nil)
				os.Exit(1)
			}

			fieldName := field.Names[0].Name
			builder.WriteString("\tif node." + fieldName + " != nil {\n")
			builder.WriteString("\t\tfor _, item := range node." + fieldName + " {\n")
			// builder.WriteString("\t\t\treuse.Free[" + eltIdent.Name + "](item, nil)\n")
			builder.WriteString("\t\t\treuse.Free[" + eltTypeName + "](item, nil)\n")
			builder.WriteString("\t\t}\n")
			builder.WriteString("\t}\n")

		} else if starExpr, ok := field.Type.(*ast.StarExpr); ok {

			eltTypeName, ok := getElementType(starExpr)
			if !ok {
				// 无法处理这个元素类型, 打印它并且退出
				ast.Fprint(os.Stdout, nil, starExpr, nil)
				os.Exit(1)
			}

			// 对于类型为 *ast.StarExpr (即指针类型)，我们只需要一个调用来释放指针指向的对象。
			fieldName := field.Names[0].Name
			builder.WriteString("\tif node." + fieldName + " != nil {\n")
			builder.WriteString("\t\treuse.Free[" + eltTypeName + "](node." + fieldName + ", nil)\n")
			builder.WriteString("\t}\n")
		}
	}

	builder.WriteString("}\n\n")

	return builder.String()
}

// hasFormatMethod 检查类型是否有 Format 方法
func hasFormatMethod(file *ast.File, structName string) bool {
	for _, decl := range file.Decls {
		if funcDecl, ok := decl.(*ast.FuncDecl); ok {
			// 检查函数是方法，有接收者，并且接收者类型与我们的结构体名称匹配
			if funcDecl.Recv != nil && len(funcDecl.Recv.List) > 0 {
				recvType := funcDecl.Recv.List[0].Type
				// 当接收者是指针类型时
				if starExpr, ok := recvType.(*ast.StarExpr); ok {
					if ident, ok := starExpr.X.(*ast.Ident); ok && ident.Name == structName {
						if funcDecl.Name.Name == "Format" {
							return true
						}
					}
				}
				// 当接收者是非指针类型时
				if ident, ok := recvType.(*ast.Ident); ok && ident.Name == structName {
					if funcDecl.Name.Name == "Format" {
						return true
					}
				}
			}
		}
	}
	return false
}