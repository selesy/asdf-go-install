SHELL := bash
.ONESHELL:
.SHELLFLAGS := -eu -o pipefail -c
.DELETE_ON_ERROR:
MAKEFLAGS += --warn-undefined-variables
MAKEFLAGS += --no-builtin-rules

clean:
	rm -rf bin/ || true
	rm -rf lib/ || true
.PHONY: clean

build: generate
	mkdir bin || true
	mkdir lib || true
	mkdir lib/commands || true
	go build -o bin/asdf-go-install .
	ln -s asdf-go-install bin/download || true
	ln -s asdf-go-install bin/install || true
	ln -s asdf-go-install bin/list-all || true
	ln -s asdf-go-install lib/commands/add || true
.PHONY: build

generate:
	go generate ./...
.PHONY: generate
