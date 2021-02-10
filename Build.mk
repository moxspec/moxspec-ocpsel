.PHONY: all bin test clean

include $(CURDIR)/Config.mk

all: bin

bin: test
	$(GO) build -o bin/ocpsel *.go

test:
	$(GO) test -coverprofile c.out ./...

clean:
	-rm -rf bin
