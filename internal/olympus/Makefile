all: olympus

olympus: version.go
	go build

version.go:
	go generate

check:
	go test -coverprofile=cover.out
	go vet

clean:
	rm -Rf cover.out

.PHONY: olympus check clean version.go
