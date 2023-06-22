all: install

build:
	@echo "Building Cosi"
	@go build -mod readonly -o build/cosi main.go

build-rocks:
	@echo "Building Cosi with RocksDB"
	@CGO_ENABLED=1 go build -mod readonly -tags rocksdb,grocksdb_clean_link -o build/cosi main.go

install:
	@echo "Installing Cosi"
	@go install -mod readonly ./...

clean:
	rm -rf build

.PHONY: all isntall clean build
