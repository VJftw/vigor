MAKEFILE_DIR:=$(dir $(realpath $(firstword $(MAKEFILE_LIST))))


.PHONY: serve main.wasm

wasm_exec.js:
	cp $(shell go env GOROOT)/misc/wasm/wasm_exec.js $@

main.wasm:
	GOARCH=wasm GOOS=js go build -o main.wasm

build: wasm_exec.js main.wasm
	rm -rf $@
	mkdir $@
	cp index.html $@/index.html
	cp wasm_exec.js $@/wasm_exec.js
	cp main.wasm $@/main.wasm

serve: build
	go run github.com/m3ng9i/ran@v0.1.6 \
		-gzip=true \
		-bind-ip=127.0.0.1 \
		-root $(realpath $(MAKEFILE_DIR)/build)
