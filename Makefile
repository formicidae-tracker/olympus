all: olympus olympus-stream-notifier/olympus-stream-notifier check webapp

clean:
	rm -f olympus

olympus-stream-notifier/olympus-stream-notifier: olympus-stream-notifier/*.go
	cd olympus-stream-notifier && go build

olympus: *.go olympuspb/*.go olympuspb/*.proto
	go generate
	go build

check:
	go test
	go vet
	make -C olympuspb check

webapp:
	make -C webapp

.PHONY: clean webapp
