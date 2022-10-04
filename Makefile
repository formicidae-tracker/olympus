all: olympus webapp olympus-stream-notifier/olympus-stream-notifier check

clean:
	rm -f olympus

olympus-stream-notifier/olympus-stream-notifier: olympus-stream-notifier/*.go
	cd olympus-stream-notifier && go build

olympus: *.go
	go generate
	go build

check:
	go test
	go vet

webapp:
	make -C webapp

.PHONY: clean webapp
