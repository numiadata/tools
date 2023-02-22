all: install

build:
	@echo "Building Cosi"
	@go build  -tags=pebbledb -mod readonly $(BUILD_FLAGS) -o build/cosi main.go -tags=pebbledb

install:
	@echo "Installing Cosi"
	@go install -tags=pebbledb -mod readonly ./... 

clean:
	rm -rf build

.PHONY: all isntall clean build
