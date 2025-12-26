_ := $(shell mkdir -p .make bin)
export GOWORK := off

DEVCTL    ?= go tool devctl
GINKGO    ?= go tool ginkgo
GOMOD2NIX ?= go tool gomod2nix
NIX       ?= nix

MODULES := docker github protofs

# GO_SRC != $(DEVCTL) list --go
GO_SRC != find . -type f -path '*.go'

ifeq ($(CI),)
TEST_FLAGS := --label-filter !E2E
else
TEST_FLAGS := --github-output --race --trace --coverprofile=cover.profile
endif

build:
	$(NIX) build .#aferox .#aferox-docker .#aferox-github .#aferox-protofs

test: .make/test
tidy: go.sum ${MODULES:%=%/go.sum}
deps: gomod2nix.toml ${MODULES:%=%/gomod2nix.toml}

test_all:
	$(GINKGO) run -r ./

%/go.sum: %/go.mod ${GO_SRC}
	go -C $* mod tidy

go.sum: go.mod ${GO_SRC}
	go mod tidy

.PHONY: gomod2nix.toml ${MODULES:%=%/gomod2nix.toml}
gomod2nix.toml ${MODULES:%=%/gomod2nix.toml}:
	$(GOMOD2NIX) generate --dir ${@D}

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

.make/test: $(filter-out ${MODULES:%=./%/%},${GO_SRC})
	$(GINKGO) run ${TEST_FLAGS} $(sort $(dir $?))
	@touch $@
