package internal

import (
	"github.com/iancoleman/strcase"
	"github.com/thecodedproject/gopkg"
	"github.com/thecodedproject/gopkg/tmpl"

	"github.com/thecodedproject/servicegen/internal/proto"
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
					{ Import: "errors" },
					{ Import: "flag" },
					{ Import: "time" },
					{ Import: "google.golang.org/grpc/connectivity" },
				},
				Vars: []gopkg.DeclVar{
					{
						Name: "address",
						Type: gopkg.TypeUnnamedLiteral{},
						LiteralValue: `flag.String("` + s.Name + `_grpc_address", "", "host:port of ` + s.Name + ` gRPC service")`,
					},
				},
				Types: []gopkg.DeclType{
					{
						Name: "grpcClient",
						Type: gopkg.TypeStruct{
							Fields: []gopkg.DeclVar{
								{
									Name: "rpcConn",
									Type: gopkg.TypePointer{
										ValueType: gopkg.TypeNamed{
											Name: "ClientConn",
											Import: "google.golang.org/grpc",
										},
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

	pbNewClientFunc := s.PbImport.Alias + ".New" + strcase.ToCamel(s.Name) + "Client"

	clientFuncs := make([]gopkg.DeclFunc, 0, len(s.ApiFuncs)+2)
	clientFuncs = append(
		clientFuncs,
		gopkg.DeclFunc{
			Name: "New",
			ReturnArgs: tmpl.UnnamedReturnArgs(
				gopkg.TypePointer{
					ValueType: gopkg.TypeNamed{
						Name: "grpcClient",
					},
				},
				gopkg.TypeError{},
			),
			BodyTmpl: `
	conn, err := grpc.Dial(*address, grpc.WithInsecure())
	if err != nil {
		{{FuncReturnDefaultsWithErr}}
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	for {
		if conn.GetState() == connectivity.Ready {
			break
		}
		if !conn.WaitForStateChange(ctx, conn.GetState()) {
			err := errors.New("grpc timeout whilst connecting")
			{{FuncReturnDefaultsWithErr}}
		}
	}

	return &grpcClient{
		rpcConn: conn,
		rpcClient: ` + pbNewClientFunc + `(conn),
	}, nil
`,
		},
		gopkg.DeclFunc{
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
				{
					Name: "conn",
					Type: gopkg.TypePointer{
						ValueType: gopkg.TypeNamed{
							Name: "ClientConn",
							Import: "google.golang.org/grpc",
						},
					},
				},
			},
			ReturnArgs: tmpl.UnnamedReturnArgs(
				gopkg.TypePointer{
					ValueType: gopkg.TypeNamed{
						Name: "grpcClient",
					},
				},
			),
			BodyTmpl: `
	return &grpcClient{
		rpcConn: conn,
		rpcClient: ` + pbNewClientFunc + `(conn),
	}
`,
		},
	)

	for _, f := range s.ApiFuncs {

		clientVarName := "_c"
		f.Receiver.VarName = clientVarName
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
			Req proto.Message
			Resp proto.Message
		}{
			PbAlias: s.PbImport.Alias,
			Req: reqMessage,
			Resp: respMessage,
		}
		f.BodyTmpl = `
{{- $pbAlias := .BodyData.PbAlias}}
	res, err := ` + clientVarName + `.rpcClient.{{.Name}}(
		ctx,
		&{{.BodyData.PbAlias}}.{{.BodyData.Req.Name}}{
{{- range .BodyData.Req.Fields}}
	{{- if .IsNestedMessage}}
			{{ToCamel .Name}}: {{$pbAlias}}.{{ToCamel .Type}}ToProto({{ToLowerCamel .Name}}),
	{{- else}}
			{{ToCamel .Name}}: {{ToLowerCamel .Name}},
	{{- end}}
{{- end}}
		},
	)
	if err != nil {
		{{FuncReturnDefaultsWithErr}}
	}

	return {{range .BodyData.Resp.Fields}}
{{- if .IsNestedMessage -}}
	{{$pbAlias}}.{{ToCamel .Type}}FromProto(res.{{ToCamel .Name}}), {{else -}}
	res.{{ToCamel .Name}}, {{end -}}
{{- end}}nil
`

		clientFuncs = append(clientFuncs, f)
	}

	return clientFuncs, nil
}
