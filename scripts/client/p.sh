#!/usr/bin/env bash

python parse.py ../protocol/
python parse_protobuf.py ../protocol
python parse_type_convert.py ../protocol
go fmt api.go
go fmt proto.go
go fmt dispatcher.go
go fmt errcode.go
go fmt const.go
go fmt decode.go
go fmt *grpc_api.go
go fmt *_type_convert.go
mv api.go ../src/protocol/
mv proto.go ../src/protocol/
mv errcode.go ../src/protocol/
mv decode.go ../src/protocol/
mv dispatcher.go ../src/gs/
mv const.go ../src/protocol/
mv *grpc_api.go ../src/grpc_api/
mv *_type_convert.go ../src/protobuf
