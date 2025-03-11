_ := $(shell mkdir -p .make bin)

WORKING_DIR := $(shell pwd)
LOCALBIN    := ${WORKING_DIR}/bin

DEVCTL := go tool devctl
GINKGO := go tool ginkgo

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

go.sum: go.mod $(shell $(DEVCTL) list --go)
	go mod tidy

%_suite_test.go:
	cd $(dir $@) && $(GINKGO) bootstrap

%_test.go:
	cd $(dir $@) && $(GINKGO) generate $(notdir $*)

.envrc: hack/example.envrc
	cp $< $@

.make/build: $(shell $(DEVCTL) list --go --exclude-tests)
	go build ./...
	@touch $@

.make/test: $(shell $(DEVCTL) list --go)
	$(GINKGO) run ${TEST_FLAGS} $(sort $(dir $?))
	@touch $@
