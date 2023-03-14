package internal

import (
	"path"

	"github.com/iancoleman/strcase"
	"github.com/thecodedproject/gopkg"
	"github.com/thecodedproject/gopkg/tmpl"
)

func fileClientTestFiles(
	s serviceDefinition,
) func() ([]gopkg.FileContents, error) {


	return func() ([]gopkg.FileContents, error) {
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
			s.ResourcesImport,
			s.PbImport,
			{
				Import: path.Join(s.ImportPath, "client", "local"),
				Alias: "client_local",
			},
			{
				Import: path.Join(s.ImportPath, "client", "grpc"),
				Alias: "client_grpc",
			},
			{
				Import: path.Join(s.ImportPath, "server"),
			},
			{ Import: "time" },
			{ Import: "context" },
			{ Import: "google.golang.org/grpc" },
			{ Import: "google.golang.org/grpc/connectivity" },
			{ Import: "github.com/stretchr/testify/require" },
			{ Import: "net" },
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
			{
				Name: "TestClientGRPCSuite",
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
				Name: "SetupTest",
				Receiver: gopkg.FuncReceiver{"ts", "TestClientGRPCSuite", true},
				// TODO: insert imports + type names dynamically in this func body
				BodyTmpl: `
	ts.createClient = func(r resources.Resources) basic.Client {
		return setupGRPCClient(ts.T(), r)
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
			{
				Name: "TestClientGRPC",
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
	suite.Run(t, new(TestClientGRPCSuite))
`,
			},
			{
				Name: "setupGRPCClient",
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
					s.ResourcesDecl,
				},
				ReturnArgs: tmpl.UnnamedReturnArgs(
					gopkg.TypeNamed{
						Name: "Client",
						Import: s.ImportPath,
						ValueType: gopkg.TypeInterface{},
					},
				),
				BodyTmpl: `
	serverAddr := setupGRPCServer(t, r)
	conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	for {
		if conn.GetState() == connectivity.Ready {
			break
		}

		if !conn.WaitForStateChange(ctx, conn.GetState()) {
			require.Fail(t, "grpc timeout whilst connecting")
		}
	}

	client := client_grpc.NewForTesting(t, conn)
	return client
`,
			},
			{
				Name: "setupGRPCServer",
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
					s.ResourcesDecl,
				},
				ReturnArgs: tmpl.UnnamedReturnArgs(
					gopkg.TypeString{},
				),
				BodyTmpl: `
	listener, err := net.Listen("tcp", "localhost:0")
	require.NoError(t, err)

	grpcSrv := grpc.NewServer()
	t.Cleanup(grpcSrv.GracefulStop)

	service := server.New(r)
	` + s.PbImport.Alias + `.Register` + strcase.ToCamel(s.Name) + `Server(grpcSrv, service)

	go func() {
		err := grpcSrv.Serve(listener)
		require.NoError(t, err)
	}()

	return listener.Addr().String()
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
				//c := ts.createClient(
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
