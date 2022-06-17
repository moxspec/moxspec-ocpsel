#!/usr/bin/make -f
.PHONY: help bin test clean

GO := go

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' Makefile | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-10s\033[0m %s\n", $$1, $$2}'

bin: ## build the binary
	$(GO) build -o bin/ocpsel *.go

test: ## test
	$(GO) test -coverprofile c.out ./...

clean: ## clean all artifacts
	-rm -rf bin/
	-rm -rf c.out
