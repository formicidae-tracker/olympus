all: olympus olympus-stream-notifier/olympus-stream-notifier check

clean:
	rm -f olympus

olympus-stream-notifier/olympus-stream-notifier: olympus-stream-notifier/*.go
	cd olympus-stream-notifier && go build

olympus: *.go api/*.go api/*.proto api/examples/*.go
	go generate
	go build

check:
	go test -coverprofile=cover.out
	go vet
	make -C api check

webapp:
	make -C webapp

.PHONY: clean webapp
