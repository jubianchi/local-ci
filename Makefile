GO := go

GO_BUILDFLAGS := -v -ldflags "-s -w"
GO_TESTFLAGS := -v -cover
GO_GETFLAGS := -v
GO_VETFLAGS := -v -all -source -shadow=true -shadowstrict

GOFMT := gofmt
GOFMTFLAGS := -s -w -l

GOPATH := $(shell pwd)
PKGPATH := $(GOPATH)/src/local-ci
BINPATH := $(GOPATH)/bin

SRCS = $(PKGPATH)/main.go $(wildcard $(PKGPATH)/**/*.go)

.PHONY: verify
verify: fmt vet

.PHONY: build build/alpine build/darwin build/linux build/windows
build: build/alpine build/darwin build/linux build/windows
build/alpine: $(BINPATH)/alpine/local-ci
build/darwin: $(BINPATH)/darwin/local-ci
build/linux: $(BINPATH)/linux/local-ci
build/windows: $(BINPATH)/windows/local-ci

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
