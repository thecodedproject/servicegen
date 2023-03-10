package internal

import (
	"errors"

	"github.com/thecodedproject/gopkg"
	"github.com/thecodedproject/gopkg/tmpl"
)

func fileResources(
	s serviceDefinition,
) func() ([]gopkg.FileContents, error) {

	return func() ([]gopkg.FileContents, error) {

		resourcedDeclType, err := declTypeFromTypeNamed(s.ResourcesDecl.Type)
		if err != nil {
			return nil, err
		}

		return []gopkg.FileContents{
			{
				Filepath: "resources/resources.go",
				PackageName: "resources",
				Types: []gopkg.DeclType{
					{
						Name: "resources",
						Type: gopkg.TypeStruct{
						},
					},
				},
				Functions: []gopkg.DeclFunc{
					{
						Name: "New",
						ReturnArgs: tmpl.UnnamedReturnArgs(
							gopkg.TypePointer{
								ValueType: gopkg.TypeNamed{
									Name: "resources",
									ValueType: gopkg.TypeStruct{},
								},
							},
						),
						BodyTmpl: `
	return &resources{}
`,
					},
				},
			},
			{
				Filepath: "resources/resources_impl.go",
				PackageName: "resources",
				Types: []gopkg.DeclType{
					resourcedDeclType,
				},
				Functions: []gopkg.DeclFunc{
					// TODO: Make NewForTesting take in a struct of all the resource fields
					{
						Name: "NewForTesting",
						Args: []gopkg.DeclVar{
							{
								Name: "_",
								Type: gopkg.TypePointer{
									ValueType: gopkg.TypeNamed{
										Name: "T",
										Import: "testing",
									},
								},
							},
						},
						ReturnArgs: tmpl.UnnamedReturnArgs(
							gopkg.TypePointer{
								ValueType: gopkg.TypeNamed{
									Name: "resources",
									ValueType: gopkg.TypeStruct{},
								},
							},
						),
						BodyTmpl: `
	return &resources{}
`,
					},
				},
			},
		}, nil
	}
}

func declTypeFromTypeNamed(t gopkg.Type) (gopkg.DeclType, error) {

	tNamed, ok := t.(gopkg.TypeNamed)
	if !ok {
		return gopkg.DeclType{}, errors.New("declTypeFromTypeNamed: Expected gopkg.TypeNamed")
	}

	return gopkg.DeclType{
		Name: tNamed.Name,
		Import: tNamed.Import,
		Type: tNamed.ValueType,
	}, nil
}
