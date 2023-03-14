package internal

import (
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
		types = append(types, gopkg.DeclType{
			Name: m.Name,
			Type: gopkg.TypeStruct{},
		})
	}

	return types, nil
}

