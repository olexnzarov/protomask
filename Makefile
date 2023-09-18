build-proto:
	rm -rf ./internal/pbtest/*.pb.go
	protoc ./internal/pbtest/*.proto \
		--go_out=. \
		--go-grpc_out=. \
		--go_opt=paths=source_relative \
		--go-grpc_opt=paths=source_relative 