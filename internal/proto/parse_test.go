package proto_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/thecodedproject/servicegen/internal/proto"
)

func TestParse(t *testing.T) {

	testCases := []struct{
		Name string
		ProtoFile string
		Expected proto.Service
	}{
		{
			Name: "Two different methods",
			ProtoFile: "./testdata/TestParse_two_different_methods.proto",
			Expected: proto.Service{
				ProtoPackage: "examplepb",
				ServiceName: "Some",
				Methods: []proto.Method{
					{
						Name: "Ping",
						RequestMessage: "PingRequest",
						ResponseMessage: "PingResponse",
					},
					{
						Name: "Pong",
						RequestMessage: "PongRequest",
						ResponseMessage: "PongResponse",
					},
				},
				Messages: []proto.Message{
					{
						Name: "PingRequest",
						Fields: []proto.Field{
							{
								Name: "int64_value",
								Type: "int64",
							},
							{
								Name: "string_value",
								Type: "string",
							},
						},
					},
					{
						Name: "PingResponse",
					},
					{
						Name: "PongRequest",
					},
					{
						Name: "PongResponse",
						Fields: []proto.Field{
							{
								Name: "int64_value",
								Type: "int64",
							},
							{
								Name: "string_value",
								Type: "string",
							},
						},
					},
				},
			},
		},
		{
			Name: "Nested messages",
			ProtoFile: "./testdata/TestParse_nested_messages.proto",
			Expected: proto.Service{
				ProtoPackage: "nestedpb",
				ServiceName: "NestedService",
				Methods: []proto.Method{
					{
						Name: "Ping",
						RequestMessage: "PingRequest",
						ResponseMessage: "PingResponse",
					},
				},
				Messages: []proto.Message{
					{
						Name: "PingRequest",
						Fields: []proto.Field{
							{
								Name: "some_nested_value",
								Type: "NestedVal",
								IsNestedMessage: true,
							},
						},
					},
					{
						Name: "PingResponse",
						Fields: []proto.Field{
							{
								Name: "some_other_value",
								Type: "OtherNestedVal",
								IsNestedMessage: true,
							},
						},
					},
					{
						Name: "NestedVal",
						Fields: []proto.Field{
							{
								Name: "some_value",
								Type: "int64",
							},
						},
					},
					{
						Name: "OtherNestedVal",
						Fields: []proto.Field{
							{
								Name: "some_string",
								Type: "string",
							},
						},
					},
				},
			},
		},
	}

	for _, test := range testCases {
		t.Run(test.Name, func(t *testing.T) {

			i, err := proto.Parse(test.ProtoFile)
			require.NoError(t, err)

			assert.Equal(t, test.Expected, i)
		})
	}
}

