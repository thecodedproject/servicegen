package internal

import (
	"github.com/iancoleman/strcase"
	"github.com/thecodedproject/gopkg"
	"github.com/thecodedproject/gopkg/tmpl"

	"github.com/thecodedproject/servicegen/internal/proto"
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
						Name: "grpcServer",
						Type: gopkg.TypeStruct{
							Embeds: []gopkg.Type{
								gopkg.TypeNamed{
									Name: "Unimplemented" + strcase.ToCamel(s.Name) + "Server",
									Import: s.PbImport.Import,
								},
							},
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

	serverFuncs := make([]gopkg.DeclFunc, 0, len(s.ApiFuncs)+1)

	serverFuncs = append(serverFuncs, gopkg.DeclFunc{
		Name: "New",
		Args: []gopkg.DeclVar{
			s.ResourcesDecl,
		},
		ReturnArgs: tmpl.UnnamedReturnArgs(
			gopkg.TypePointer{
				ValueType: gopkg.TypeNamed{
					Name: "grpcServer",
				},
			},
		),
		BodyTmpl: `
	return &grpcServer{
		r: r,
	}
`,
	})

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

		serverVarName := "_s"
		f.Receiver.VarName = serverVarName
		f.Receiver.TypeName = "grpcServer"
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
			Req proto.Message
			Resp proto.Message
		}{
			PbAlias: s.PbImport.Alias,
			Req: reqMessage,
			Resp: respMessage,
		}
		f.BodyTmpl = `
{{- $pbAlias := .BodyData.PbAlias}}
	{{range .BodyData.Resp.Fields}}{{ToLowerCamel .Name}}, {{end}}err := internal.{{.Name}}(
		ctx,
		` + serverVarName + `.r,
{{- range .BodyData.Req.Fields}}
	{{- if .IsNestedMessage}}
		{{$pbAlias}}.{{ToCamel .Type}}FromProto(req.{{ToCamel .Name}}),
	{{- else}}
		req.{{ToCamel .Name}},
	{{- end}}
{{- end}}
	)
	if err != nil {
		return nil, err
	}

	return &{{.BodyData.PbAlias}}.{{.BodyData.Resp.Name}}{
{{- range .BodyData.Resp.Fields}}
	{{- if .IsNestedMessage}}
		{{ToCamel .Name}}: {{$pbAlias}}.{{ToCamel .Type}}ToProto({{ToLowerCamel .Name}}),
	{{- else }}
		{{ToCamel .Name}}: {{ToLowerCamel .Name}},
	{{- end}}
{{- end}}
	}, nil
`

		serverFuncs = append(serverFuncs, f)
	}

	return serverFuncs, nil
}
