GO ?= go

PKGS = ./
BIN = ./cmd

all: fmt install test

fmt:
	$(GO) fmt $(PKGS)

test:
	$(GO) test -race $(PKGS)

install: build mvbin

mvbin: 
	mv spritesheet $(GOPATH)/bin/

build:
	$(GO) build -i -o spritesheet $(BIN)
	
