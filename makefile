.PHONY: .FORCE
include env.mk

GO?=go
PWD=$(shell pwd)
GOBIN := $(GOPATH)/bin
GOPATH := $(PWD)

VERSION ?= M3
REVISION ?= $(REVISION)
BUILD_TIME=$(shell date '+%Y-%m-%d_%H:%M:%S_%Z')
VER_PKG=server/base/option

# go linker flags, possible options see https://pkg.go.dev/cmd/link
LDFLAGS="\
-X $(VER_PKG)._VERSION_=$(VERSION) \
-X $(VER_PKG)._REV_=$(REVISION) \
-X $(VER_PKG)._BuildTime=$(BUILD_TIME) "

# go compiler flags, possible options see https://pkg.go.dev/cmd/compile
GCFLAGS ?= ""

GOFLAGS ?= -gcflags=$(GCFLAGS) -ldflags=$(LDFLAGS)

dev:
	go get github.com/golang/protobuf/protoc-gen-go
	go get github.com/gogo/protobuf/protoc-gen-gofast

genpbgo:
	clang-format -i $(ALL_PB_SRC)
	mkdir -p $(PB_GO_SRC_DIR)
	protoc --proto_path=$(PB_SRC_DIR) --gofast_out=$(PB_GO_SRC_DIR) $(ALL_PB_SRC)

genmsg:
	python ./scripts/client/gen_msg_go.py ${ALL_PB_SRC_OPT} --out $(PB_GO_SRC_DIR)/msg.go
	@if [ ! -e bin/genmsgid ]; then  GOPATH=$(GOPATH) $(GO) install Doudou/app/genmsgid ; fi
	bin/genmsgid -path=$(PB_GO_SRC_DIR)/msgid.def