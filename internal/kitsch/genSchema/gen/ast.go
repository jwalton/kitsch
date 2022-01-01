package gen

import (
	"go/ast"
	"go/token"
)

// findStruct finds the struct with the given name in the given file and returns it.
func findStruct(file *ast.File, structName string) *ast.StructType {
	var result *ast.StructType

	ast.Inspect(file, func(n ast.Node) bool {
		// Find our struct declaration
		declaration, ok := n.(*ast.GenDecl)
		if !ok || declaration.Tok != token.TYPE {
			return true
		}

		typeSpec, ok := declaration.Specs[0].(*ast.TypeSpec)
		if !ok || typeSpec.Name.Name != structName {
			return true
		}

		structType, ok := typeSpec.Type.(*ast.StructType)
		if !ok {
			return true
		}

		result = structType
		return false
	})

	return result
}
