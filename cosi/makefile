all: install

build:
	@echo "Building Cosi"
	@go build -mod readonly $(BUILD_FLAGS) -o build/cosi main.go

install:
	@echo "Installing Cosi"
	@go install -mod readonly ./...

clean:
	rm -rf build

.PHONY: all isntall clean build
