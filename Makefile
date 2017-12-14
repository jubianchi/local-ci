GO := go

GO_BUILDFLAGS ?= -v -ldflags "-s -w"
GO_TESTFLAGS ?= -v -cover
GO_GETFLAGS ?= -v
GO_VETFLAGS ?= -v -all -source -shadow=true -shadowstrict

GOFMT := gofmt
GOFMTFLAGS := -s -w -l

GOPATH := $(shell pwd)
PKGPATH := $(GOPATH)/src/local-ci
BINPATH := $(GOPATH)/bin

SRCS = $(PKGPATH)/main.go $(wildcard $(PKGPATH)/**/*.go)

.PHONY: verify
verify: fmt vet

.PHONY: build build/alpine build/darwin build/linux
build: build/alpine build/darwin build/linux build/windows
build/alpine: $(BINPATH)/alpine/local-ci
build/darwin: $(BINPATH)/darwin/local-ci
build/linux: $(BINPATH)/linux/local-ci
build/windows: $(BINPATH)/windows/local-ci
	mv $< $<.exe

.PHONY: clean
clean:
	rm -rf bin/ pkg/ src/github.com gopkk.in/

.PHONY: fmt
fmt:
	@$(GOFMT) $(GOFMTFLAGS) $(PKGPATH)

.PHONY: vet
vet:
	GOPATH=$(GOPATH) $(GO) vet $(GO_VETFLAGS) ./src/local-ci/...

.PHONY: test
test:
	GOPATH=$(GOPATH) $(GO) test $(GO_TESTFLAGS) local-ci/...

.PHONY: vendor
vendor:
	GOPATH=$(GOPATH) $(GO) get $(GO_GETFLAGS) ./src/local-ci/...

$(BINPATH)/%/local-ci: $(SRCS) vendor
	mkdir -p $(dir $@)
	GOPATH=$(GOPATH) GOOS=$(shell echo $* | sed s/alpine/linux/) $(GO) build $(GO_BUILDFLAGS) -o $@ $<

CHANGELOG.md: $(BINPATH)/clog .clog.toml
	$< --setversion 1.0.0-beta.1

$(BINPATH)/clog:
	mkdir -p $(BINPATH)
ifeq ($(shell uname),Darwin)
	wget -q -O $@.tar.gz https://github.com/clog-tool/clog-cli/releases/download/v0.9.3/clog-v0.9.3-$(shell uname -m)-apple-darwin.tar.gz
else
	wget -q -O $@.tar.gz https://github.com/clog-tool/clog-cli/releases/download/v0.9.3/clog-v0.9.3-$(shell uname -m)-unknown-linux-gnu.tar.gz
endif
	tar xzvf $@.tar.gz -C bin/
	rm -f $@.tar.gz
