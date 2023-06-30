all: olympus_service.pb.go olympus_service_grpc.pb.go

%.pb.go: olympus_service.proto
	go generate

.PHONY: check

check:
	go test
	go vet
