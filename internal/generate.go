package internal

import (
	"errors"
	"flag"

	"github.com/iancoleman/strcase"
	"github.com/thecodedproject/gopkg"
	"github.com/thecodedproject/gopkg/tmpl"

	"github.com/thecodedproject/servicegen/internal/proto"
)

var (
	protoFile = flag.String("proto", "", "path to proto file")
)

func Generate() error {

	flag.Parse()

	if *protoFile == "" {
		return errors.New("`proto` must be set")
	}

	apiProto, err := proto.Parse(*protoFile)
	if err != nil {
		return err
	}

	s, err := createServiceDefinition(apiProto)
	if err != nil {
		return err
	}


	return gopkg.LintAndGenerate([]gopkg.FileContents{
		makeApiFile(s),
		makeClientLocalClientFile(s),
	})
}

type serviceDefinition struct {
	Name string
	ImportPath string
	ApiFuncs []gopkg.DeclFunc
}

func createServiceDefinition(
	apiProto proto.Service,
) (serviceDefinition, error) {

	importPath, err := gopkg.PackageImportPath(".")
	if err != nil {
		return serviceDefinition{}, err
	}

	apiFuncs, err := apiFuncSignatures(apiProto)
	if err != nil {
		return serviceDefinition{}, err
	}

	return serviceDefinition{
		Name: strcase.ToSnake(apiProto.ServiceName),
		ImportPath: importPath,
		ApiFuncs: apiFuncs,
	}, nil
}

func apiFuncSignatures(
	apiProto proto.Service,
) ([]gopkg.DeclFunc, error) {

	funcs := make([]gopkg.DeclFunc, len(apiProto.Methods))

	for iM, m := range apiProto.Methods {

		reqMessage, err := apiProto.MethodRequestMessage(m.Name)
		if err != nil {
			return nil, err
		}

		respMessage, err := apiProto.MethodResponseMessage(m.Name)
		if err != nil {
			return nil, err
		}

		funcs[iM] = tmpl.FuncWithContextAndError(
			m.Name,
			argsFromProtoMessage(reqMessage),
			argsFromProtoMessage(respMessage),
		)
	}

	return funcs, nil
}

func argsFromProtoMessage(
	m proto.Message,
) []gopkg.DeclVar {

	args := make([]gopkg.DeclVar, len(m.Fields))

	for iF, f := range m.Fields {
		args[iF] = gopkg.DeclVar{
			Name: strcase.ToLowerCamel(f.Name),
			Type: goTypeFromProtoType(f.Type),
		}
	}

	return args
}

func goTypeFromProtoType(
	protoType string,
) gopkg.Type {

	switch protoType {
	case "bool": return gopkg.TypeBool{}
	case "bytes":
		return gopkg.TypeArray{
			ValueType: gopkg.TypeByte{},
		}
	case "float": return gopkg.TypeFloat64{}
	case "int32": return gopkg.TypeInt32{}
	case "int64": return gopkg.TypeInt64{}
	case "string": return gopkg.TypeString{}
	}

	return gopkg.TypeUnknownNamed{
		Name: protoType,
	}
}

func makeApiFile(
	s serviceDefinition,
) gopkg.FileContents {

	return gopkg.FileContents{
		Filepath: "api.go",
		PackageName: s.Name,
		Types: []gopkg.DeclType{
			{
				Name: strcase.ToCamel(s.Name),
				Type: gopkg.TypeInterface{
					Funcs: s.ApiFuncs,
				},
			},
		},
	}
}

