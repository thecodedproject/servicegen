package internal

import (
	"github.com/iancoleman/strcase"
	"github.com/thecodedproject/gopkg"
)

func fileApi(
	s serviceDefinition,
) gopkg.FileContents {

	return gopkg.FileContents{
		Filepath: "api.go",
		PackageName: s.Name,
		Types: []gopkg.DeclType{
			{
				Name: strcase.ToCamel(s.Name),
				Type: gopkg.TypeInterface{
					Funcs: s.ApiFuncs,
				},
			},
		},
	}
}
