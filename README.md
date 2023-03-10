# servicegen
Golang generator for creating services

# dependencies

* `protoc` generator must be installed

# protoc generated files
Generated with:
```
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative path/to.proto
```
