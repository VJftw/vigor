.PHONY: test wasmbrowsertest

GOBIN := $(shell go env GOBIN)
GOPATH := $(shell go env GOPATH)
ifndef GOBIN
	GOBIN=$(GOPATH)/bin
endif
PATH := $(GOBIN):$(PATH)

wasmbrowsertest: export GOOS = $(shell go env GOOS)
wasmbrowsertest: export GOARCH = $(shell go env GOARCH)
wasmbrowsertest:
ifeq (,$(wildcard $(GOBIN)/go_js_wasm_exec))
	go install github.com/agnivade/wasmbrowsertest@latest
	mv $(GOBIN)/wasmbrowsertest $(GOBIN)/go_js_wasm_exec
endif


test: wasmbrowsertest
test: export GOOS = js
test: export GOARCH = wasm
test: PACKAGES = $(shell go list ../...)
test:
	@rm -rf .coverage
	@mkdir -p .coverage
	@echo "mode: count" > .coverage/merged.cover.out

	@$(foreach pkg,$(PACKAGES), \
		mkdir -p .coverage/$(pkg)/ \
		&& go test -v -cover -covermode=count -coverpkg=$(pkg) -coverprofile=.coverage/$(pkg)/cover.out $(pkg)/tests || exit 1 \
	;)

	@$(foreach pkg,$(PACKAGES), \
		cat .coverage/$(pkg)/cover.out | grep -v "mode:" >> .coverage/merged.cover.out \
	;)


	@go tool cover -html=.coverage/merged.cover.out -o .coverage/coverage.html
	go tool cover -func=.coverage/merged.cover.out
