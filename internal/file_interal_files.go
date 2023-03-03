package internal

import (
	"path/filepath"

	"github.com/iancoleman/strcase"
	"github.com/thecodedproject/gopkg"
)

func fileInternalFiles(
	s serviceDefinition,
) []gopkg.FileContents {

	internalFiles := make([]gopkg.FileContents, 0, len(s.ApiFuncs))
	for _, f := range s.ApiFuncs {

		f.BodyTmpl = `
	{{FuncReturnDefaults}}
`

		internalFiles = append(internalFiles, gopkg.FileContents{
			Filepath: filepath.Join("internal", strcase.ToSnake(f.Name) + ".go"),
			PackageName: s.Name,
			Functions: []gopkg.DeclFunc{
				f,
			},
		})
	}

	return internalFiles
}
