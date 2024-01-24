package reuseutils

import (
	"bytes"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
)

func FormatFile(filepath string) error {
	fset := token.NewFileSet()
	// 解析Go代码文件
	file, err := parser.ParseFile(fset, filepath, nil, parser.ParseComments)
	if err != nil {
		return err
	}

	// cmap 创建一个节点到注释的映射
	cmap := ast.NewCommentMap(fset, file, file.Comments)

	// 方法映射：类型名称到方法列表
	typeMethods := make(map[string][]*ast.FuncDecl)

	// 其他声明列表（不包括方法）
	var otherDecls []ast.Decl
	for _, decl := range file.Decls {
		switch tDecl := decl.(type) {
		case *ast.GenDecl:
			otherDecls = append(otherDecls, tDecl)
		case *ast.FuncDecl:
			if tDecl.Recv != nil && len(tDecl.Recv.List) > 0 {
				recvType := tDecl.Recv.List[0].Type
				if starExpr, isStar := recvType.(*ast.StarExpr); isStar {
					recvType = starExpr.X
				}
				if ident, isIdent := recvType.(*ast.Ident); isIdent {
					typeName := ident.Name
					typeMethods[typeName] = append(typeMethods[typeName], tDecl)
					continue
				}
			}
			otherDecls = append(otherDecls, tDecl)
		}
	}

	// 重建文件声明列表，将类型声明和对应的方法按正确顺序组合在一起
	var newDecls []ast.Decl
	for _, decl := range otherDecls {
		newDecls = append(newDecls, decl)
		if genDecl, ok := decl.(*ast.GenDecl); ok {
			for _, spec := range genDecl.Specs {
				if typeSpec, ok := spec.(*ast.TypeSpec); ok {
					if methods, ok := typeMethods[typeSpec.Name.Name]; ok {
						for _, m := range methods {
							newDecls = append(newDecls, ast.Decl(m))
						}
					}
				}
			}
		}
	}

	// 更新文件的声明列表和注释
	file.Decls = newDecls
	file.Comments = cmap.Filter(file).Comments()

	// 格式化和输出更新后的代码到相同文件
	var buf bytes.Buffer
	if err := format.Node(&buf, fset, file); err != nil {
		return err
	}
	return os.WriteFile(filepath, buf.Bytes(), 0o644)
}

// func Format(fp string, fi os.FileInfo, err error) error {
// 	if err != nil || fi.IsDir() || strings.HasSuffix(fp, "_test.go") {
// 		return err
// 	}
// 	if strings.HasSuffix(fi.Name(), ".go") {
// 		log.Printf("Processing file: %s", fp)
// 		err := reuseutils.FormatFile(fp)
// 		if err != nil {
// 			log.Printf("Error processing file '%s': %v", fp, err)
// 			return err
// 		}
// 	}
// 	return nil
// }

// func main() {
// 	root := "/Users/jensen/Projects/matrixorigin/matrixone/pkg/sql/parsers/tree/alter.go" // 或者可以通过命令行参数传入
// 	err := filepath.Walk(root, Format)
// 	if err != nil {
// 		log.Fatalf("Error walking through the root directory: %v", err)
// 	}
// }
