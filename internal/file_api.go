package internal

import (
	"github.com/thecodedproject/gopkg"
)

func fileApi(
	s serviceDefinition,
) func() ([]gopkg.FileContents, error) {

	return func() ([]gopkg.FileContents, error) {
		return []gopkg.FileContents{
			{
				Filepath: "api.go",
				PackageName: s.Name,
				PackageImportPath: s.ImportPath,
				Types: []gopkg.DeclType{
					{
						Name: "Client",
						Type: gopkg.TypeInterface{
							Funcs: s.ApiFuncs,
						},
					},
				},
			},
		}, nil
	}
}

