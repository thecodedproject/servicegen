package internal

import (
	"path"

	"github.com/thecodedproject/gopkg"
	"github.com/thecodedproject/gopkg/tmpl"
)

func fileClientLocalClient(
	s serviceDefinition,
) func() ([]gopkg.FileContents, error) {

	return func() ([]gopkg.FileContents, error) {
		internalImportPath := path.Join(s.ImportPath, "internal")
		resourcesImportPath := path.Join(s.ImportPath, "resources")

		funcs := []gopkg.DeclFunc{
			{
				Name: "New",
				Args: []gopkg.DeclVar{
					s.ResourcesDecl,
				},
				ReturnArgs: tmpl.UnnamedReturnArgs(
					gopkg.TypePointer{
						ValueType: gopkg.TypeNamed{
							Name: "client",
						},
					},
				),
				BodyTmpl: `
	return &client{
		r: r,
	}
`,
			},
		}

		clientVarName := "_c"

		for _, f := range s.ApiFuncs {
			f.Receiver = gopkg.FuncReceiver{
				VarName: clientVarName,
				TypeName: "client",
				IsPointer: true,
			}

			f.BodyData = internalFuncCallParams(f.Args, clientVarName)
			f.BodyTmpl = `
	return internal.{{.Name}}(
	{{- range .BodyData}}
		{{.}},
	{{- end}}
	)
`
			funcs = append(funcs, f)
		}

		return []gopkg.FileContents{
			{
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
			},
		}, nil
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
	clientVarName string,
) []string {

	params := make([]string, len(apiFuncArgs))
	for i, arg := range apiFuncArgs {
		params[i] = arg.Name
	}

	if len(params) < 2 {
		return append(params, clientVarName + ".r")
	}

	params = append(params[:2], params[1:]...)
	params[1] = clientVarName + ".r"
	return params
}
