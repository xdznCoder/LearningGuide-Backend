protoc -I . post.proto --go_out=. --go-grpc_out=.
protoc -I . user.proto --go_out=. --go-grpc_out=.