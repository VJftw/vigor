MAKEFILE_DIR:=$(dir $(realpath $(firstword $(MAKEFILE_LIST))))

.PHONY: build serve test

build:
	@rm -rf $@
	@mkdir $@
	@cp $(dir $(shell go env GOWORK))/examples/_common/index.html $@/index.html
	@cp $(shell go env GOROOT)/misc/wasm/wasm_exec.js $@/wasm_exec.js
	@GOARCH=wasm GOOS=js go build -o $@/main.wasm

test: build
	go test -v ./test/.

serve: build
	go run github.com/m3ng9i/ran@v0.1.6 \
		-gzip=true \
		-bind-ip=127.0.0.1 \
		-root $(realpath $(MAKEFILE_DIR)/build)
