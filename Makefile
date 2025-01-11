build_proto:
	protoc --go_out=./internal/notification --go_opt=paths=source_relative \
	--go-grpc_out=./internal/notification --go-grpc_opt=paths=source_relative \
	proto/notifications.proto