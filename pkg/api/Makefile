all: olympus_service.pb.go olympus_service_grpc.pb.go mock_api/olympus_server.go

mock_api/olympus_server.go: olympus_service.proto
	go generate

%.pb.go: olympus_service.proto
	go generate

.PHONY: check

check: *.go
	go test
	go vet
