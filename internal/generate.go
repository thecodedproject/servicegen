package internal

import (
	"errors"
	"flag"
	"path"

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

	var files []gopkg.FileContents
	files, err = appendFileContents(
		files,
		fileApi(s),
		fileClientGrpcClient(s),
		fileClientLocalClient(s),
		fileClientTestFiles(s),
		fileInternalFiles(s),
		filePbTypeConversion(s),
		fileResources(s),
		fileServerGrpcServer(s),
		fileTypes(s),
	)
	if err != nil {
		return err
	}

	return gopkg.LintAndGenerate(files)
}

func appendFileContents(
	files []gopkg.FileContents,
	fileFuncs ...func() ([]gopkg.FileContents, error),
) ([]gopkg.FileContents, error) {

	for _, fileFunc := range fileFuncs {
		newFiles, err := fileFunc()
		if err != nil {
			return nil, err
		}

		files = append(files, newFiles...)
	}
	return files, nil
}

type serviceDefinition struct {
	Name string
	ImportPath string
	ApiFuncs []gopkg.DeclFunc
	ResourcesDecl gopkg.DeclVar
	ApiProto proto.Service

	ClientImport gopkg.ImportAndAlias
	InternalImport gopkg.ImportAndAlias
	PbImport gopkg.ImportAndAlias
	ResourcesImport gopkg.ImportAndAlias
}

func createServiceDefinition(
	apiProto proto.Service,
) (serviceDefinition, error) {

	importPath, err := gopkg.PackageImportPath(".")
	if err != nil {
		return serviceDefinition{}, err
	}

	apiFuncs, err := apiFuncSignatures(apiProto, importPath)
	if err != nil {
		return serviceDefinition{}, err
	}

	resImp := makeImportWithAlias(importPath, "resources")

	sName := strcase.ToSnake(apiProto.ServiceName)

	return serviceDefinition{
		Name: sName,
		ImportPath: importPath,
		ApiFuncs: apiFuncs,
		ResourcesDecl: gopkg.DeclVar{
			Name: "r",
			Type: gopkg.TypeNamed{
				Name: "Resources",
				Import: resImp.Import,
				ValueType: gopkg.TypeInterface{},
			},
		},
		ApiProto: apiProto,
		ClientImport: makeImportWithAlias(importPath, "client"),
		InternalImport: makeImportWithAlias(importPath, "internal"),
		PbImport: makeImportWithAlias(importPath, "pb"),
		ResourcesImport: resImp,
	}, nil
}

func apiFuncSignatures(
	apiProto proto.Service,
	serviceImportPath string,
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
			argsFromProtoMessage(reqMessage, serviceImportPath),
			unnamedArgsFromProtoMessage(respMessage, serviceImportPath),
		)
	}

	return funcs, nil
}

func makeImportWithAlias(
	baseImportPath string,
	subPath string,
) gopkg.ImportAndAlias {

	i := path.Join(baseImportPath, subPath)
	return gopkg.ImportAndAlias{
		Import: i,
		Alias: subPath,
	}
}

func argsFromProtoMessage(
	m proto.Message,
	serviceImportPath string,
) []gopkg.DeclVar {

	args := make([]gopkg.DeclVar, len(m.Fields))

	for iF, f := range m.Fields {
		args[iF] = gopkg.DeclVar{
			Name: strcase.ToLowerCamel(f.Name),
			Type: goTypeFromProtoType(f.Type, serviceImportPath),
		}
	}

	return args
}

func unnamedArgsFromProtoMessage(
	m proto.Message,
	serviceImportPath string,
) []gopkg.DeclVar {

	args := make([]gopkg.DeclVar, len(m.Fields))

	for iF, f := range m.Fields {
		args[iF] = gopkg.DeclVar{
			Type: goTypeFromProtoType(f.Type, serviceImportPath),
		}
	}

	return args
}

func goTypeFromProtoType(
	protoType string,
	serviceImportPath string,
) gopkg.Type {

	switch protoType {
	case "bool": return gopkg.TypeBool{}
	case "bytes":
		return gopkg.TypeArray{
			ValueType: gopkg.TypeByte{},
		}
	case "float": return gopkg.TypeFloat32{}
	case "int32": return gopkg.TypeInt32{}
	case "int64": return gopkg.TypeInt64{}
	case "string": return gopkg.TypeString{}
	}

	return gopkg.TypeNamed{
		Name: protoType,
		Import: serviceImportPath,
		ValueType: gopkg.TypeStruct{},
	}
}

