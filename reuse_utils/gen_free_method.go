
package reuseutils

import (
	"go/ast"
	"strings"
)


func GenerateFree(_ *ast.StructType, structName string) string {
	var builder strings.Builder
	builder.WriteString("func (node *" + structName + ") Free() {\n")
	builder.WriteString("	reuse.Free[" + structName + "](node, nil)\n}\n")
	return builder.String()
}