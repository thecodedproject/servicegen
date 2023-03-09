package internal

import (
	"path"
	"path/filepath"

	"github.com/iancoleman/strcase"
	"github.com/thecodedproject/gopkg"
)

func fileInternalFiles(
	s serviceDefinition,
) func() ([]gopkg.FileContents, error) {

	return func() ([]gopkg.FileContents, error) {
		resourcesImportPath := path.Join(s.ImportPath, "resources")

		resourceType := gopkg.DeclVar{
			Name: "r",
			Type: gopkg.TypeNamed{
				Name: "Resources",
				Import: resourcesImportPath,
				ValueType: gopkg.TypeStruct{},
			},
		}

		internalFiles := make([]gopkg.FileContents, 0, len(s.ApiFuncs))
		for _, f := range s.ApiFuncs {

			f.BodyTmpl = `
	{{FuncReturnDefaults}}
`

			// Add `Resources` as the second arg in this function signature
			if len(f.Args) < 2 {
				f.Args = append(f.Args, resourceType)
			} else {
				f.Args = append(f.Args[:2], f.Args[1:]...)
				f.Args[1] = resourceType
			}

			internalFiles = append(internalFiles, gopkg.FileContents{
				Filepath: filepath.Join("internal", strcase.ToSnake(f.Name) + ".go"),
				PackageName: s.Name,
				Functions: []gopkg.DeclFunc{
					f,
				},
			})
		}

		return internalFiles, nil
	}
}
