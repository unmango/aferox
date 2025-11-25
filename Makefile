_ := $(shell mkdir -p .make bin)
export GOWORK := off

GO        ?= go
DEVCTL    ?= $(GO) tool devctl
GINKGO    ?= $(GO) tool ginkgo
GOMOD2NIX ?= $(GO) tool gomod2nix

MODULES := docker github protofs

# GO_SRC != $(DEVCTL) list --go
GO_SRC != find . -type f -path '*.go'

ifeq ($(CI),)
TEST_FLAGS := --label-filter !E2E
else
TEST_FLAGS := --github-output --race --trace --coverprofile=cover.profile
endif

build: .make/build
test: .make/test
tidy: go.sum ${MODULES:%=%/go.sum}

test_all:
	$(GINKGO) run -r ./

%/go.sum: %/go.mod ${GO_SRC}
	go -C $* mod tidy

gomod2nix.toml: go.mod go.sum
	$(GOMOD2NIX)
docker/gomod2nix.toml: docker/go.mod docker/go.sum
	$(GOMOD2NIX) --dir ${@D}
github/gomod2nix.toml: github/go.mod github/go.sum
	$(GOMOD2NIX) --dir ${@D}
protofs/gomod2nix.toml: protofs/go.mod protofs/go.sum
	$(GOMOD2NIX) --dir ${@D}

go.sum: go.mod ${GO_SRC}
	go mod tidy

go.work: export GOWORK :=
go.work: ${MODULES:%=%/go.mod}
	go work init
	go work use . ${MODULES}
go.work.sum: go.work
	go work sync

%_suite_test.go:
	cd $(dir $@) && $(GINKGO) bootstrap
%_test.go:
	cd $(dir $@) && $(GINKGO) generate $(notdir $*)

.envrc: hack/example.envrc
	cp $< $@

.make/build: $(filter-out %_test.go,${GO_SRC})
	go build ./...
	@touch $@

.make/test: $(filter-out ${MODULES:%=./%/%},${GO_SRC})
	$(GINKGO) run ${TEST_FLAGS} $(sort $(dir $?))
	@touch $@
