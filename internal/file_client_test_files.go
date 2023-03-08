package internal

import (
	"path"

	"github.com/iancoleman/strcase"
	"github.com/thecodedproject/gopkg"
	"github.com/thecodedproject/gopkg/tmpl"
)

func fileClientTestFiles(
	s serviceDefinition,
) ([]gopkg.FileContents, error) {

	files := []gopkg.FileContents{
		fileClientTestCommon(s),
	}

	testFiles, err := fileClientTestForAPIFuncs(s)
	if err != nil {
		return nil, err
	}

	files = append(files, testFiles...)

	return files, nil
}

func fileClientTestCommon(
	s serviceDefinition,
) gopkg.FileContents {

	return gopkg.FileContents{
		Filepath: "client/client_test.go",
		PackageName: "client_test",
		Imports: []gopkg.ImportAndAlias{
			{
				Import: s.ImportPath,
				Alias: s.Name,
			},
			{
				Import: s.ResourcesImport,
				Alias: "resources",
			},
			{
				Import: path.Join(s.ImportPath, "client", "local"),
				Alias: "client_local",
			},
		},
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
									s.ResourcesDecl.Type,
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
	ts.createClient = func(r resources.Resources) basic.Client {
		return client_local.New(r)
	}
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
) ([]gopkg.FileContents, error) {

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
				{
					Import: "github.com/stretchr/testify/require",
					Alias: "require",
				},
			},
			Functions: []gopkg.DeclFunc{
				{
					Name: "Test" + strcase.ToCamel(f.Name),
					Receiver: gopkg.FuncReceiver{"ts", "clientSuite", true},
					BodyData: f,
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

				require.Fail(t, "TODO: Implement test...")
				//c = ts.createClient(
				//	resources.NewForTesting(t),
				//)

				//ctx := context.Background()
				//..., err := c.{{.BodyData.Name}}(ctx, ...)
				//require.NoError(t, err)
			})
		}
	`,
				},
			},
		}
	}

	return files, nil
}
