GOPATH:=$(shell pwd)
GO:=go
GOFLAGS:=-v -p 1
APPLICATION:=alien_invasion
TARGETS:=bin/${APPLICATION}

default: ${TARGETS}

all: clean default

version/version.go: .version
	sh -c 'mkdir -p version && ./version.sh .version >$@'

bin/${APPLICATION}: version/version.go
	@echo "========== Compiling artifacts for: $@ =========="
	sh -c '$(GO) build $(GOFLAGS) -o $@ github.com/hartsp2000/${APPLICATION}'

format:
	@echo "========== Formatting artifacts =========="
	@for f in `find . -name \*.go`; do echo "Reformatting $$f"; go fmt $$f || { echo "Issue found in Golang code. Please fix them and try again."; exit 1; }; done

vet:
	@GO_MODULE_NAME=`head -1 go.mod | awk '{ print $$2 }'`; for path_to_test in `find . -type d -print`; do  { ls >/dev/null 2>&1 -al $${path_to_test}/*.go || continue ; };  path_to_test=`echo $${path_to_test} | sed 's/^.\///g'`; go_path=$${GO_MODULE_NAME}/$${path_to_test}; echo -n "Vetting $${go_path} ... "; ${GO} vet $${go_path}; if [ $$? -ne 0 ]; then echo \"Issue found in Golang code. Please fix them and try again.\"; exit 2; else echo "[ OK ]"; fi done

clean:
	@echo "========== Cleaning artifacts =========="
	@echo "Deleting generated binary files ..."; for binary in ${TARGETS}; do if [ -f "$${binary}" ]; then rm -f $${binary} && echo $${binary}; fi; done
	@echo "Deleting generated version files ..."; for version_dir in /version ; do if [ -d "version" ]; then rm -Rf version && echo version; fi; done
	echo "Deleting backup files: "
	find . -name \*~ -exec rm -f {}  \;
