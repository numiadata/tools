#!/usr/bin/make -f

all: install

build:
	@echo "Building erebus..."
	@go build -mod readonly $(BUILD_FLAGS) -o build/erebus main.go

install:
	@echo "Installing erebus..."
	@go install -mod readonly ./...

clean:
	rm -rf build

.PHONY: all build install
