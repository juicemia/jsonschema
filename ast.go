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
		processStruct(typ, fset)
	}

}

func processStruct(st *ast.StructType, fset *token.FileSet) {
	schema := Type{
		Type:       "object",
		Properties: map[string]*Type{},
	}

	fields := st.Fields
	for _, fld := range fields.List {
		t := getNodeContent(fld.Type, fset)

		ftyp := Type{
			Type: t,
		}
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

func getNodeContent(node ast.Node, fset *token.FileSet) string {
	begin := fset.PositionFor(node.Pos(), true)
	end := fset.PositionFor(node.End(), true)
	spew.Dump(begin, end)

	tmpf, err := os.Open(begin.Filename)
	if err != nil {
		panic(err)
	}

	buf := make([]byte, end.Offset-begin.Offset)
	_, err = tmpf.ReadAt(buf, int64(begin.Offset))
	if err != nil {
		panic(err)
	}

	return string(buf)
}
