_ := $(shell mkdir -p .make bin)

WORKING_DIR := $(shell pwd)
LOCALBIN    := ${WORKING_DIR}/bin

export GOBIN := ${LOCALBIN}

DEVCTL := ${LOCALBIN}/devctl
GINKGO := ${LOCALBIN}/ginkgo

ifeq ($(CI),)
TEST_FLAGS := --label-filter !E2E
else
TEST_FLAGS := --github-output --race --trace --coverprofile=cover.profile
endif

build: .make/build
test: .make/test
tidy: go.sum

test_all:
	$(GINKGO) run -r ./

go.sum: go.mod $(shell $(DEVCTL) list --go) | bin/devctl
	go mod tidy

%_suite_test.go: | bin/ginkgo
	cd $(dir $@) && $(GINKGO) bootstrap

%_test.go: | bin/ginkgo
	cd $(dir $@) && $(GINKGO) generate $(notdir $*)

bin/ginkgo: go.mod
	go install github.com/onsi/ginkgo/v2/ginkgo

bin/devctl: .versions/devctl
	go install github.com/unmango/devctl/cmd@v$(shell cat $<)
	mv ${LOCALBIN}/cmd $@

.envrc: hack/example.envrc
	cp $< $@

.make/build: $(shell $(DEVCTL) list --go --exclude-tests) | bin/devctl
	go build ./...
	@touch $@

.make/test: $(shell $(DEVCTL) list --go) | bin/ginkgo bin/devctl
	$(GINKGO) run ${TEST_FLAGS} $(sort $(dir $?))
	@touch $@
