package codegen

import (
	"log"
	"strings"

	"github.com/dave/jennifer/jen"
	"github.com/vektah/gqlparser"
	"github.com/vektah/gqlparser/ast"
)

type query struct {
	*ast.QueryDocument
	string
}

func ConvertToGo(schema string, queries ...string) *jen.File {
	s := gqlparser.MustLoadSchema(&ast.Source{Input: schema})
	qs := make([]query, len(queries))
	for i, q := range queries {
		qs[i] = query{gqlparser.MustLoadQuery(s, q), q}
	}
	graphqlt := "github.com/flexzuu/graphqlt"
	f := jen.NewFile("graphql")
	f.Type().Id("Client").Struct(
		jen.Op("*").Qual(graphqlt, "Client"),
	)
	f.Comment("NewClient makes a new Client capable of making GraphQL requests.")
	f.Func().Id("NewClient").Params(jen.Id("endpoint").String(), jen.Id("opts").Op("...").Qual(graphqlt, "ClientOption")).Op("*").Id("Client").Block(
		jen.Id("c").Op(":=").Op("&").Id("Client").Values(jen.Qual(graphqlt, "NewClient").Call(jen.Id("endpoint"), jen.Id("opts").Op("..."))),
		jen.Return().Id("c"),
	)
	for _, q := range qs {
		for _, o := range q.Operations {
			name := strings.Title(o.Name)

			responseName := name + "Response"
			f.Type().Id(responseName).Struct(convertSelectionSet(o.SelectionSet, "Query", s)...)
			f.Func().Params(
				jen.Id("c").Op("*").Id("Client"),
			).Id(name).Params(
				jen.Id("ctx").Qual("context", "Context"),
				jen.Id("c").String(),
			).Params(jen.Id(responseName), jen.Error()).Block(
				jen.Id("req").Op(":=").Qual(graphqlt, "NewRequest").Call(jen.Lit(q.string)),
				jen.Var().Id("respData").Id(responseName),
				jen.Id("err").Op(":=").Id("c").Dot("Run").Call(jen.Id("ctx"), jen.Id("req"), jen.Op("&").Id("respData")),
				jen.Return(jen.Id("respData"), jen.Id("err")),
			)
			log.Println(o.Name)
		}
	}
	return f
}

func convertSelectionSet(ss ast.SelectionSet, parentName string, schema *ast.Schema) []jen.Code {
	var parent *ast.Definition
	if parentName == "Query" {
		parent = schema.Query
	} else {
		p := schema.Types[parentName]

		parent = p
	}
	if parent == nil {
		panic("not found")
	}
	fs := make([]jen.Code, 0)
	for _, sel := range ss {

		switch s := sel.(type) {
		case *ast.Field:
			f := parent.Fields.ForName(s.Name)
			t := schema.Types[f.Type.Name()]
			if t == nil {
				panic("not found")
			}
			if !t.IsLeafType() {
				fs = append(fs, jen.Id(strings.Title(s.Name)).Struct(convertSelectionSet(s.SelectionSet, t.Name, schema)...))
			} else {
				fs = append(fs, jen.Id(strings.Title(s.Name)).Add(toGoType(t.Name)))
			}
		default:
			panic("unkown type")
		}
	}
	return fs
}

func toGoType(t string) jen.Code {
	switch t {
	case "String":
		return jen.String()
	default:
		panic("unkown type")
	}

}
