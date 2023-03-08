package resources

import (
	"testing"
)

type Resources interface{}

func NewForTesting(_ *testing.T) *resources {

	return &resources{
	}
}
