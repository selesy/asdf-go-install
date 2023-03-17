SHELL := bash
.ONESHELL:
.SHELLFLAGS := -eu -o pipefail -c
.DELETE_ON_ERROR:
MAKEFLAGS += --warn-undefined-variables
MAKEFLAGS += --no-builtin-rules

clean:
	rm -rf bin/ || true
.PHONY: clean

build: generate
	mkdir bin || true
	go build -o bin/asdf-go-install .
	ln -s asdf-go-install bin/download || true
	ln -s asdf-go-install bin/install || true
	ln -s asdf-go-install bin/list-all || true
.PHONY: build

generate:
	go generate ./...
.PHONY: generate
