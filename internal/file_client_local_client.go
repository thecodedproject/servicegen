package internal

import (
	"path"

	"github.com/thecodedproject/gopkg"
)

func fileClientLocalClient(
	s serviceDefinition,
) gopkg.FileContents {

	internalImportPath := path.Join(s.ImportPath, "internal")
	resourcesImportPath := path.Join(s.ImportPath, "resources")

	funcs := make([]gopkg.DeclFunc, 0, len(s.ApiFuncs))
	for _, f := range s.ApiFuncs {
		f.Receiver = gopkg.FuncReceiver{
			VarName: "c",
			TypeName: "client",
			IsPointer: true,
		}

		f.BodyData = internalFuncCallParams(f.Args)
		f.BodyTmpl = `
	return internal.{{.Name}}(
	{{- range .BodyData}}
		{{.}},
	{{- end}}
	)
`
		funcs = append(funcs, f)
	}

	return gopkg.FileContents{
		Filepath: "client/local/client.go",
		PackageName: "local",
		Imports: []gopkg.ImportAndAlias{
			{
				Alias: "internal",
				Import: internalImportPath,
			},
		},
		Types: []gopkg.DeclType{
			{
				Name: "client",
				Type: gopkg.TypeStruct{
					Fields: []gopkg.DeclVar{
						{
							Name: "r",
							Type: gopkg.TypeNamed{
								Name: "Resources",
								Import: resourcesImportPath,
							},
						},
					},
				},
			},
		},
		Functions: funcs,
	}
}

// internalFuncCallParams returns the parameter names needed to call
// the internal function from the client function.
//
// The internal function will have the same parameters as ther api function
// with addition of the service's resources struct at position two.
// i.e. for an API func:
//		Method(ctx, one, two)
// the internal call will be
//		internal.Method(ctx, c.r, one two)
// where `c.r` is the resources struct coming from the function receiver
// `c *Client`
func internalFuncCallParams(
	apiFuncArgs []gopkg.DeclVar,
) []string {

	params := make([]string, len(apiFuncArgs))
	for i, arg := range apiFuncArgs {
		params[i] = arg.Name
	}

	if len(params) < 2 {
		return append(params, "c.r")
	}

	params = append(params[:2], params[1:]...)
	params[1] = "c.r"
	return params
}
