// Utility functions to help working with gopkg
package internal

import (
	"errors"

	"github.com/thecodedproject/gopkg"
)

func declTypeFromTypeNamed(t gopkg.Type) (gopkg.DeclType, error) {

	tNamed, ok := t.(gopkg.TypeNamed)
	if !ok {
		return gopkg.DeclType{}, errors.New("declTypeFromTypeNamed: Expected gopkg.TypeNamed")
	}

	return gopkg.DeclType{
		Name: tNamed.Name,
		Import: tNamed.Import,
		Type: tNamed.ValueType,
	}, nil
}
