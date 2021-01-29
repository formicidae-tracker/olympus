VERSION := $(shell git describe)
LDFLAGS :=-ldflags "-X 'main.OLYMPUS_VERSION=$(VERSION)'"

all: olympus webapp olympus-stream-notifier/olympus-stream-notifier

clean:
	rm -f olympus

olympus-stream-notifier/olympus-stream-notifier: olympus-stream-notifier/*.go
	cd olympus-stream-notifier && go build

olympus: *.go
	go build $(LDFLAGS)

webapp:
	make -C webapp

.PHONY: clean webapp
