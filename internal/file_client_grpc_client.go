package internal

import (
	"github.com/iancoleman/strcase"
	"github.com/thecodedproject/gopkg"
)

func fileClientGrpcClient(
	s serviceDefinition,
) func() ([]gopkg.FileContents, error) {

	return func() ([]gopkg.FileContents, error) {

		funcs, err := makeGrpcClientFuncs(s)
		if err != nil {
			return nil, err
		}

		return []gopkg.FileContents{
			{
				Filepath: "client/grpc/grpc_client.go",
				PackageName: "grpc",
				Imports: []gopkg.ImportAndAlias{
					s.PbImport,
				},
				Types: []gopkg.DeclType{
					{
						Name: "grpcClient",
						Type: gopkg.TypeStruct{
							Fields: []gopkg.DeclVar{
								{
									Name: "rpcConn",
									Type: gopkg.TypeNamed{
										Name: "ClientConn",
										Import: "google.golang.org/grpc",
									},
								},
								{
									Name: "rpcClient",
									Type: gopkg.TypeNamed{
										Name: strcase.ToCamel(s.Name) + "Client",
										Import: s.PbImport.Import,
									},
								},
							},
						},
					},
				},
				Functions: funcs,
			},
		}, nil
	}
}

func makeGrpcClientFuncs(
	s serviceDefinition,
) ([]gopkg.DeclFunc, error) {

	clientFuncs := make([]gopkg.DeclFunc, 0, len(s.ApiFuncs))
	for _, f := range s.ApiFuncs {

		f.Receiver.VarName = "c"
		f.Receiver.TypeName = "grpcClient"
		f.Receiver.IsPointer = true

		reqMessage, err := s.ApiProto.MethodRequestMessage(f.Name)
		if err != nil {
			return nil, err
		}

		respMessage, err := s.ApiProto.MethodResponseMessage(f.Name)
		if err != nil {
			return nil, err
		}

		f.BodyData = struct{
			PbAlias string
			ReqName string
			ReqArgNames []string
			RespArgNames []string
		}{
			PbAlias: s.PbImport.Alias,
			ReqName: reqMessage.Name,
			ReqArgNames: reqMessage.FieldNames(),
			RespArgNames: respMessage.FieldNames(),
		}
		f.BodyTmpl = `


	res, err := c.rpcClient.{{.Name}}(
		ctx,
		&{{.BodyData.PbAlias}}.{{.BodyData.ReqName}}{
{{- range .BodyData.ReqArgNames}}
			{{ToCamel .}}: {{ToLowerCamel .}},
{{- end}}
		},
	)
	if err != nil {
		{{FuncReturnDefaultsWithErr}}
	}

	return {{range .BodyData.RespArgNames}}res.{{ToCamel .}}, {{end}}nil
`

		clientFuncs = append(clientFuncs, f)
	}

	return clientFuncs, nil
}
