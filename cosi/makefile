


all: install



build:
	@echo "Building Cosi"
	@go build -mod readonly $(BUILD_FLAGS) -o build/cosi main.go

install:
	@echo "Installing Cosi"
	@go install -mod readonly $(BUILD_FLAGS) ./...

clean:
	rm -rf build

.PHONY: all lint test race msan tools clean build
