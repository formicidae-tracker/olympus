all: olympus check

clean:
	rm -f cmd/olympus/olympus

olympus:
	make -C pkg/api
	make -C internal/olympus
	make -C cmd/olympus

check:
	make -C pkg/api check
	make -C internal/olympus check

webapp:
	make -C webapp

.PHONY: olympus clean webapp
