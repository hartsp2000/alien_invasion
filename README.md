::::::::::::::::::::
:: Alien Invasion ::
::::::::::::::::::::

Golang simulator for a theoretical alien invasion

Notes:
------
For the ease of portability in this demonstration package, golang modules have NOT been used.

FAQ for potential compilation errors (due to platform differences):
-------------------------------------------------------------------
If there is a compiliation error while running make regarding modules, please update the local golang
config by running the following command:

$ go env -w GO111MODULE=auto

If there is a compliation error: go/pkg/mod/golang.org/x/sys@(etc..etc..etc)/unix/zsyscall_darwin_amd64.go:##:##: //go:linkname must refer to declared function or variable, run the following command to update the package:

$ go get -u golang.org/x/sys

If there is a compilation error: $GOPATH/go.mod exists but should not, unset the GOPATH:

$ unset GOPATH


Cleaning:
---------
$ make clean

Building/Compiling:
-------------------
$ make

Formatting:
-----------
$ make format

Vetting:
--------
$ make vet

