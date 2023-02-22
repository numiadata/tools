all: install

build:
	@echo "Building Cosi"
	@go build -mod readonly $(BUILD_FLAGS) -o build/cosi main.go -tags=pebbledb

install:
	@echo "Installing Cosi"
	@go install -mod readonly ./... -tags=pebbledb

clean:
	rm -rf build

.PHONY: all isntall clean build
