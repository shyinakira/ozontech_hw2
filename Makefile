run-all:
	(cd products && go run ./cmd/products/main.go) & \
	(cd cart && go run ./cmd/cart/main.go)

generate:
	protoc --proto_path=products \
		--go_out=products/gen/go --go_opt=paths=source_relative \
		--go-grpc_out=products/gen/go --go-grpc_opt=paths=source_relative \
		products/proto/product_service.proto