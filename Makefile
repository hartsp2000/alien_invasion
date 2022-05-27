GOPATH:=$(shell pwd)
APPLICATION:=alien_invasion
GO:=go
GOFLAGS:=-v -p 1
TARGETS:=bin/${APPLICATION}

default: ${TARGETS}

all: clean default

src/version/version.go: .version
	sh -c 'mkdir -p version && ./bin/version.sh .version >$@'

format:
	@echo "========== Formatting artifacts =========="
	@for f in `find . -name \*.go`; do echo "Reformatting $$f"; go fmt $$f || { echo "Issue found in Golang code. Please fix them and try again."; exit 1; }; done

vet:
	@echo "========== Vetting Golang artifacts =========="
	@for file_to_test in `find . -name *.go -print`; do { printf "Vetting $${file_to_test}..."; $(GO) vet $${file_to_test}; if [ $$? -eq 0 ]; then echo "[ OK ]"; fi } done

bin/${APPLICATION}: src/version/version.go
	@echo "========== Compiling artifacts for: $@ =========="
	export GOPATH=$(shell pwd); $(GO) install $(GOFLAGS) saga.xyz/${APPLICATION}

clean:
	@echo "========== Cleaning artifacts =========="
	@echo "Deleting generated binary files ..."; for binary in ${TARGETS} ; do if [ -f "$${binary}" ]; then rm -f $${binary} && echo $${binary}; fi; done
	echo "Deleting backup files: "
	find . -name \*~ -exec rm -f {}  \;
