package internal

import (
	"github.com/iancoleman/strcase"
	"github.com/thecodedproject/gopkg"
	"github.com/thecodedproject/gopkg/tmpl"
)

func filePbTypeConversion(
	s serviceDefinition,
) func() ([]gopkg.FileContents, error) {

	return func() ([]gopkg.FileContents, error) {

		serviceImportAlias := s.Name

		funcs, err := makeTypeConversionFuncs(s, serviceImportAlias)
		if err != nil {
			return nil, err
		}

		if len(funcs) == 0 {
			return nil, nil
		}

		return []gopkg.FileContents{
			{
				Filepath: "pb/type_conversion.go",
				PackageName: s.PbImport.Alias,
				PackageImportPath: s.PbImport.Import,
				Imports: []gopkg.ImportAndAlias{
					{
						Import: s.ImportPath,
						Alias: serviceImportAlias,
					},
				},
				Functions: funcs,
			},
		}, nil
	}
}

func makeTypeConversionFuncs(
	s serviceDefinition,
	serviceImportAlias string,
) ([]gopkg.DeclFunc, error) {

	nestedMsgs := s.ApiProto.NestedMessages()

	funcs := make([]gopkg.DeclFunc, 0, len(nestedMsgs)*2)

	for _, m := range nestedMsgs {

		mName := strcase.ToCamel(m.Name)

		apiType := gopkg.TypeNamed{
			Name: mName,
			Import: s.ImportPath,
			ValueType: gopkg.TypeStruct{},
		}

		protoType := gopkg.TypePointer{
			ValueType: gopkg.TypeNamed{
				Name: mName,
				Import: s.PbImport.Import,
			},
		}

		funcs = append(
			funcs,
			gopkg.DeclFunc{
				Name: mName + "FromProto",
				Args: []gopkg.DeclVar{
					{
						Name: "v",
						Type: protoType,
					},
				},
				ReturnArgs: tmpl.UnnamedReturnArgs(
					apiType,
				),
				BodyData: m,
				BodyTmpl: `
	if v == nil {
		{{FuncReturnDefaults}}
	}

	return ` + serviceImportAlias + `.` + mName + `{
{{- range .BodyData.Fields}}
	{{- if .IsNestedMessage}}
		{{ToCamel .Name}}: {{.Type}}FromProto(v.{{ToCamel .Name}}),
	{{- else}}
		{{ToCamel .Name}}: v.{{ToCamel .Name}},
	{{- end}}
{{- end}}
	}
`,
			},
			gopkg.DeclFunc{
				Name: strcase.ToCamel(m.Name) + "ToProto",
				Args: []gopkg.DeclVar{
					{
						Name: "v",
						Type: apiType,
					},
				},
				ReturnArgs: tmpl.UnnamedReturnArgs(
					protoType,
				),
				BodyData: m,
				BodyTmpl: `
	return &` + mName + `{
{{- range .BodyData.Fields}}
	{{- if .IsNestedMessage}}
		{{ToCamel .Name}}: {{.Type}}ToProto(v.{{ToCamel .Name}}),
	{{- else}}
		{{ToCamel .Name}}: v.{{ToCamel .Name}},
	{{- end}}
{{- end}}
	}
`,
			},
		)
	}

	return funcs, nil
}

