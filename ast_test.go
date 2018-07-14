package jsonschema

import "testing"

func TestAST(t *testing.T) {
	loadTypeFromPackageFile("TestType", "github.com/juicemia/jsonschema/ast_test.go")
}

type TestType struct {
	Field string `json:"field"`
}
