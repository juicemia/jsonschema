package jsonschema

import (
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"strings"

	"github.com/davecgh/go-spew/spew"
)

func loadTypeFromPackageFile(t, pf string) {
	gopath := os.Getenv("GOPATH")
	filename := fmt.Sprintf("%v/src/%v", gopath, pf)
	fset := token.NewFileSet()
	astf, err := parser.ParseFile(fset, filename, nil, 0)
	if err != nil {
		panic(err)
	}

	ts := astf.Scope.Objects[t].Decl.(*ast.TypeSpec)
	// spew.Dump(ts)

	switch typ := ts.Type.(type) {
	case *ast.StructType:
		processStruct(typ)
	}

}

func processStruct(st *ast.StructType) {
	schema := Type{
		Type:       "object",
		Properties: map[string]*Type{},
	}

	fields := st.Fields
	for _, fld := range fields.List {
		ftyp := Type{}
		schema.Properties[parseJSONTag(fld.Tag.Value)] = &ftyp
	}

	buf, err := json.Marshal(schema)
	if err != nil {
		panic(err)
	}

	spew.Printf("%s\n", buf)
}

func parseJSONTag(tag string) string {
	spl := strings.Split(tag, " ")

	for _, t := range spl {
		if strings.HasPrefix(t, "`json") {
			return t[7 : len(t)-2]
		}
	}

	return ""
}
