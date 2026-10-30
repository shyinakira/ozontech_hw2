.PHONY: run-all run-products run-cart generate test test-race test-api

run-all:
	(cd products && go run ./cmd/products/main.go) & \
	(cd cart && go run ./cmd/cart/main.go)

run-products:
	cd products && go run ./cmd/products/main.go

run-cart:
	cd cart && go run ./cmd/cart/main.go

generate:
	protoc --proto_path=products \
		--go_out=products/gen/go --go_opt=paths=source_relative \
		--go-grpc_out=products/gen/go --go-grpc_opt=paths=source_relative \
		products/proto/product_service.proto

test:
	go test ./cart/... ./products/...

test-race:
	go test -race ./cart/... ./products/...

test-api:
	go test -v -count=1 -tags=api ./cart/test/... ./products/test/...
