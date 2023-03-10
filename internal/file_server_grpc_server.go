package internal

import (
	"path"

	"github.com/iancoleman/strcase"
	"github.com/thecodedproject/gopkg"
	"github.com/thecodedproject/gopkg/tmpl"
)

func fileServerGrpcServer(
	s serviceDefinition,
) func() ([]gopkg.FileContents, error) {

	return func() ([]gopkg.FileContents, error) {


		serverFuncs, err := makeServerFuncs(s)
		if err != nil {
			return nil, err
		}

		return []gopkg.FileContents{
			{
				Filepath: "server/grpc_server.go",
				PackageName: "server",
				Types: []gopkg.DeclType{
					{
						Name: "GRPCServer",
						Type: gopkg.TypeStruct{
							Fields: []gopkg.DeclVar{
								s.ResourcesDecl,
							},
						},
					},
				},
				Functions: serverFuncs,
			},
		}, nil
	}
}

func makeServerFuncs(
	s serviceDefinition,
) ([]gopkg.DeclFunc, error) {

	pbImport := path.Join(s.ImportPath, strcase.ToSnake(s.Name) + "pb")

	serverFuncs := make([]gopkg.DeclFunc, 0, len(s.ApiFuncs))
	for _, f := range s.ApiFuncs {

		f = tmpl.FuncWithContextAndError(
			f.Name,
			[]gopkg.DeclVar{
				{
					Name: "req",
					Type: gopkg.TypePointer{
						ValueType: gopkg.TypeNamed{
							Name: f.Name + "Request",
							Import: pbImport,
							ValueType: gopkg.TypeStruct{},
						},
					},
				},
			},
			tmpl.UnnamedReturnArgs(
				gopkg.TypePointer{
					ValueType: gopkg.TypeNamed{
						Name: f.Name + "Response",
						Import: pbImport,
						ValueType: gopkg.TypeStruct{},
					},
				},
			),
		)

		f.Receiver.VarName = "s"
		f.Receiver.TypeName = "GRPCServer"
		f.Receiver.IsPointer = true

		f.BodyTmpl = `
	{{FuncReturnDefaults}}
`

		serverFuncs = append(serverFuncs, f)
	}

	return serverFuncs, nil
}
