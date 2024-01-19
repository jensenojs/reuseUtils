package reuseutils

// import (
// 	"bytes"
// 	"go/ast"
// 	"go/parser"
// 	"go/token"
// 	"os"
// 	"path/filepath"
// 	"text/template"
// )

// var methodTemplate = template.Must(template.New("").Parse(`
// func (node {{.}}) TypeName() string {
// 	return "tree.{{.}}"
// }
// `))

// func generateMethodForType(typeName string) ([]byte, error) {
// 	var buf bytes.Buffer
// 	if err := methodTemplate.Execute(&buf, typeName); err != nil {
// 		return nil, err
// 	}
// 	return buf.Bytes(), nil
// }

// func processFile(filePath string) error {
// 	fset := token.NewFileSet()

// 	node, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
// 	if err != nil {
// 		return err
// 	}

// 	for _, f := range node.Decls {
// 		genDecl, ok := f.(*ast.GenDecl)
// 		if !ok {
// 			continue
// 		}
// 		if genDecl.Tok != token.TYPE {
// 			continue
// 		}
// 		for _, spec := range genDecl.Specs {
// 			typeSpec, ok := spec.(*ast.TypeSpec)
// 			if !ok {
// 				continue
// 			}
// 			methodCode, err := generateMethodForType(typeSpec.Name.Name)
// 			if err != nil {
// 				return err
// 			}

// 			// Append generated code to the file
// 			file, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY, 0644)
// 			if err != nil {
// 				return err
// 			}
// 			defer file.Close()

// 			if _, err := file.Write(methodCode); err != nil {
// 				return err
// 			}
// 		}
// 	}

// 	return nil
// }

// func main() {
// 	dirPath := "/Users/jensen/Projects/matrixorigin/matrixone/pkg/sql/parsers/tree"
// 	err := filepath.Walk(dirPath, func(path string, f os.FileInfo, err error) error {
// 		if err != nil {
// 			return err
// 		}
// 		if !f.IsDir() && filepath.Ext(path) == ".go" {
// 			err := processFile(path)
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
