package proto_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/thecodedproject/servicegen/internal/proto"
)

func TestService_NestedMessages(t *testing.T) {

	testCases := []struct{
		Name string
		S proto.Service
		Expected []proto.Message
	}{
		{
			Name: "empty service returns empty",
			Expected: []proto.Message{},
		},
		{
			Name: "no nested messages returns empty",
			S: proto.Service{
				Messages: []proto.Message{
					{Name: "A"},
					{Name: "B"},
					{Name: "C"},
					{Name: "D"},
				},
			},
			Expected: []proto.Message{},
		},
		{
			Name: "some nested msgs",
			S: proto.Service{
				Messages: []proto.Message{
					{
						Name: "A",
						Fields: []proto.Field{
							{
								Name: "b",
								Type: "B",
							},
							{
								Name: "non-nested",
								Type: "SomeOtherType",
							},
						},
					},
					{Name: "B"},
					{
						Name: "C",
						Fields: []proto.Field{
							{
								Name: "d",
								Type: "D",
							},
						},
					},
					{
						Name: "D",
						Fields: []proto.Field{
							{
								Name: "val",
								Type: "string",
							},
						},
					},
				},
			},
			Expected: []proto.Message{
				{Name: "B"},
				{
					Name: "D",
					Fields: []proto.Field{
						{
							Name: "val",
							Type: "string",
						},
					},
				},
			},
		},
		{
			Name: "nested msgs which are nested multiple times",
			S: proto.Service{
				Messages: []proto.Message{
					{
						Name: "A",
						Fields: []proto.Field{
							{
								Name: "b_one",
								Type: "B",
							},
							{
								Name: "b_two",
								Type: "B",
							},
							{
								Name: "non-nested",
								Type: "SomeOtherType",
							},
						},
					},
					{Name: "B"},
					{
						Name: "C",
						Fields: []proto.Field{
							{
								Name: "d",
								Type: "D",
							},
							{
								Name: "a",
								Type: "A",
							},
							{
								Name: "b",
								Type: "B",
							},
						},
					},
					{
						Name: "D",
						Fields: []proto.Field{
							{
								Name: "val",
								Type: "string",
							},
						},
					},
				},
			},
			Expected: []proto.Message{
				{
					Name: "A",
					Fields: []proto.Field{
						{
							Name: "b_one",
							Type: "B",
						},
						{
							Name: "b_two",
							Type: "B",
						},
						{
							Name: "non-nested",
							Type: "SomeOtherType",
						},
					},
				},
				{Name: "B"},
				{
					Name: "D",
					Fields: []proto.Field{
						{
							Name: "val",
							Type: "string",
						},
					},
				},
			},
		},
	}

	for _, test := range testCases {
		t.Run(test.Name, func(t *testing.T) {

			actual := test.S.NestedMessages()
			require.Equal(t, test.Expected, actual)
		})
	}
}
