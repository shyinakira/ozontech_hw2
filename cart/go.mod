module hw2/cart

go 1.27.0

replace hw2/products => ../products

require (
	github.com/stretchr/testify v1.12.1
	golang.org/x/sync v0.23.0
	google.golang.org/grpc v1.84.0
	hw2/products v0.0.0
)

require (
	go.yaml.in/yaml/v3 v3.0.5 // indirect
	golang.org/x/net v0.58.0 // indirect
	golang.org/x/sys v0.48.0 // indirect
	golang.org/x/text v0.42.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260706201446-f0a921348800 // indirect
	google.golang.org/protobuf v1.36.12 // indirect
)
