package internal

import (
	"github.com/iancoleman/strcase"
	"github.com/thecodedproject/gopkg"
)

func fileTypes(
	s serviceDefinition,
) func() ([]gopkg.FileContents, error) {

	return func() ([]gopkg.FileContents, error) {

		types, err := makeTypes(s)
		if err != nil {
			return nil, err
		}

		if len(types) == 0 {
			return nil, nil
		}

		return []gopkg.FileContents{
			{
				Filepath: "types.go",
				PackageName: s.Name,
				PackageImportPath: s.ImportPath,
				Types: types,
			},
		}, nil
	}
}

func makeTypes(
	s serviceDefinition,
) ([]gopkg.DeclType, error) {

	nestedMsgs := s.ApiProto.NestedMessages()

	types := make([]gopkg.DeclType, 0, len(nestedMsgs))

	for _, m := range nestedMsgs {

		fields := make([]gopkg.DeclVar, 0, len(m.Fields))
		for _, field := range m.Fields {
			fields = append(fields, gopkg.DeclVar{
				Name: strcase.ToCamel(field.Name),
				Type: goTypeFromProtoType(field.Type, s.ImportPath),
			})
		}

		types = append(types, gopkg.DeclType{
			Name: m.Name,
			Type: gopkg.TypeStruct{
				Fields: fields,
			},
		})
	}

	return types, nil
}

