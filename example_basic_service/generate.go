package basic

//go:generate go run ../main.go --proto basicpb/basic.proto
//go:generate protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative basicpb/basic.proto
