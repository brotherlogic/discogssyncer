protoc --proto_path ../../../ -I=./server --go_out=plugins=grpc:./server server/server.proto
