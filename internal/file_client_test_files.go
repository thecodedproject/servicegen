package internal

import (
	"path"

	"github.com/iancoleman/strcase"
	"github.com/thecodedproject/gopkg"
	"github.com/thecodedproject/gopkg/tmpl"
)

func fileClientTestFiles(
	s serviceDefinition,
) []gopkg.FileContents {

	files := []gopkg.FileContents{
		fileClientTestCommon(s),
	}

	files = append(files, fileClientTestForAPIFuncs(s)...)

	return files
}

func fileClientTestCommon(
	s serviceDefinition,
) gopkg.FileContents {

	resourcesImportPath := path.Join(s.ImportPath, "resources")

	return gopkg.FileContents{
		Filepath: "client/client_test.go",
		PackageName: "client_test",
		Types: []gopkg.DeclType{
			{
				Name: "clientSuite",
				Type: gopkg.TypeStruct{
					Embeds: []gopkg.Type{
						gopkg.TypeNamed{
							Name: "Suite",
							Import: "github.com/stretchr/testify/suite",
						},
					},
					Fields: []gopkg.DeclVar{
						{
							Name: "createClient",
							Type: gopkg.TypeFunc{
								Args: tmpl.UnnamedReturnArgs(
									gopkg.TypeNamed{
										Name: "Resources",
										Import: resourcesImportPath,
									},
								),
								ReturnArgs: tmpl.UnnamedReturnArgs(
									gopkg.TypeNamed{
										Name: "Client",
										Import: s.ImportPath,
									},
								),
							},
						},
					},
				},
			},
			{
				Name: "TestClientLocalSuite",
				Type: gopkg.TypeStruct{
					Embeds: []gopkg.Type{
						gopkg.TypeNamed{
							Name: "clientSuite",
						},
					},
				},
			},
		},
		Functions: []gopkg.DeclFunc{
			{
				Name: "SetupTest",
				Receiver: gopkg.FuncReceiver{"ts", "TestClientLocalSuite", true},
				// TODO: insert imports + type names dynamically in this func body
				BodyTmpl: `
	//ts.createClient = func(r resources.Resources) example_import_path.Client {
	//	return client_local.New(r)
	//}
`,
			},
			{
				Name: "TestClientLocal",
				Args: []gopkg.DeclVar{
					{
						Name: "t",
						Type: gopkg.TypePointer{
							ValueType: gopkg.TypeNamed{
								Name: "T",
								Import: "testing",
							},
						},
					},
				},
				BodyTmpl: `
	suite.Run(t, new(TestClientLocalSuite))
`,
			},
		},
	}
}

func fileClientTestForAPIFuncs(
	s serviceDefinition,
) []gopkg.FileContents {

	//resourcesImportPath := path.Join(s.ImportPath, "resources")

	files := make([]gopkg.FileContents, len(s.ApiFuncs))
	for i, f := range s.ApiFuncs {
		files[i] = gopkg.FileContents{
			Filepath: "client/client_" + strcase.ToSnake(f.Name) + "_test.go",
			PackageName: "client_test",
			Imports: []gopkg.ImportAndAlias{
				{
					Import: "testing",
					Alias: "testing",
				},
			},
			Functions: []gopkg.DeclFunc{
				{
					Name: "Test" + strcase.ToCamel(f.Name),
					Receiver: gopkg.FuncReceiver{"ts", "clientSuite", true},
					BodyTmpl: `
		testCases := []struct{
			Name string
		}{
			{
				Name: "empty",
			},
		}

		for _, test := range testCases {
			ts.T().Run(test.Name, func(t *testing.T) {

				return
			})
		}
	`,
				},
			},
		}
	}

	return files
}
