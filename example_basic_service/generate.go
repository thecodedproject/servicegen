package basic

//go:generate go run ../main.go --proto api.proto
//go:generate protoc --go_out=pb --go_opt=paths=source_relative --go-grpc_out=pb --go-grpc_opt=paths=source_relative api.proto
