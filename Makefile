VERSION := $(shell git describe)
LDFLAGS :=-ldflags "-X 'main.OLYMPUS_VERSION=$(VERSION)'"

all: olympus webapp

clean:
	rm -f olympus

olympus:
	go build $(LDFLAGS)

webapp:
	make -C webapp

.PHONY: clean olympus webapp
