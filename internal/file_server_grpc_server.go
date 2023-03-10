package internal

import (
	"github.com/thecodedproject/gopkg"
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

	serverFuncs := make([]gopkg.DeclFunc, 0, len(s.ApiFuncs))
	for _, f := range s.ApiFuncs {

		f.Receiver.VarName = "s"
		f.Receiver.TypeName = "GRPCServer"
		f.Receiver.IsPointer = true

		f.Args = nil

		f.ReturnArgs = nil

		serverFuncs = append(serverFuncs, f)
	}

	return serverFuncs, nil
}
