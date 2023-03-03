package internal

import (
	"path"

	"github.com/thecodedproject/gopkg"
)

func fileClientLocalClient(
	s serviceDefinition,
) gopkg.FileContents {

	internalImportPath := path.Join(s.ImportPath, "internal")

	funcs := make([]gopkg.DeclFunc, 0, len(s.ApiFuncs))
	for _, f := range s.ApiFuncs {
		f.Receiver = gopkg.FuncReceiver{
			VarName: "c",
			TypeName: "client",
			IsPointer: true,
		}

		f.BodyTmpl = `
	return internal.{{.Func.Name}}(
	{{- range .Func.Args}}
		{{.Name}},
	{{- end}}
	)
`
		funcs = append(funcs, f)
	}

	return gopkg.FileContents{
		Filepath: "client/local/client.go",
		PackageName: "local",
		Imports: []gopkg.ImportAndAlias{
			{
				Alias: "internal",
				Import: internalImportPath,
			},
		},
		Types: []gopkg.DeclType{
			{
				Name: "client",
				Type: gopkg.TypeStruct{
					Fields: []gopkg.DeclVar{
						{
							Name: "r",
							Type: gopkg.TypeUnknownNamed{
								Name: "Resource",
								Import: "my/resource",
							},
						},
					},
				},
			},
		},
		Functions: funcs,
	}
}
