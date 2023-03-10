package internal

import (
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
				Imports: []gopkg.ImportAndAlias{
					s.InternalImport,
				},
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

		f = tmpl.FuncWithContextAndError(
			f.Name,
			[]gopkg.DeclVar{
				{
					Name: "req",
					Type: gopkg.TypePointer{
						ValueType: gopkg.TypeNamed{
							Name: f.Name + "Request",
							Import: s.PbImport.Import,
							ValueType: gopkg.TypeStruct{},
						},
					},
				},
			},
			tmpl.UnnamedReturnArgs(
				gopkg.TypePointer{
					ValueType: gopkg.TypeNamed{
						Name: f.Name + "Response",
						Import: s.PbImport.Import,
						ValueType: gopkg.TypeStruct{},
					},
				},
			),
		)

		f.Receiver.VarName = "s"
		f.Receiver.TypeName = "GRPCServer"
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
			RespName string
			ReqArgNames []string
			RespArgNames []string
		}{
			PbAlias: s.PbImport.Alias,
			RespName: respMessage.Name,
			ReqArgNames: reqMessage.FieldNames(),
			RespArgNames: respMessage.FieldNames(),
		}
		f.BodyTmpl = `

	{{range .BodyData.RespArgNames}}{{ToLowerCamel .}}, {{end}}err := internal.{{.Name}}(
		ctx,
		s.r,
{{- range .BodyData.ReqArgNames}}
		req.{{ToCamel .}},
{{- end}}
	)
	if err != nil {
		return nil, err
	}

	return &{{.BodyData.PbAlias}}.{{.BodyData.RespName}}{
{{- range .BodyData.RespArgNames}}
		{{ToCamel .}}: {{ToLowerCamel .}},
{{- end}}
	}, nil
`

		serverFuncs = append(serverFuncs, f)
	}

	return serverFuncs, nil
}
