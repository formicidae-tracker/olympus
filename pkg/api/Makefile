all: olympus_service.pb.go olympus_service_grpc.pb.go

%.pb.go: olympus_service.proto generate_client_utils.go
	go generate

.PHONY: check

check: *.go
	go test
	go vet
