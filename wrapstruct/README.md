# wrapstruct
This is a simple and hardly tested CLI tool aimed to generate wrappers around
Go structs to replace public and/or private fields with similarly named getters.

__WARNING__: this tool automates creation of non-idiomatic Go code.

It can be useful in some circumstances, for example to generate an adapter
between existing interface and structure, especially, if either of them is not
under developers control, or expensive to change. Or, probably, in some other rare cases.
Otherwise, it is, probably, better to avoid use it. 

I created it mostly to experiment with go:generate and go/types package, and
put it here just in case I ever need it or I ever want to create something similar.

## Useful links on the topic:

- https://blog.gopheracademy.com/advent-2015/reducing-boilerplate-with-go-generate/
- https://github.com/golang/example/tree/master/gotypes#struct-types
- https://golang.org/pkg/go/types/
- https://godoc.org/golang.org/x/tools/go/loader#example-Config-Import
- https://github.com/ernesto-jimenez/gogen (some snippets I shamelessly borrowed there)

## Sample usage

```
# install it

go get github.com/nikolay-turpitko/x-go/wrapstruct/.../

# see usage

$GOPATH/bin/wrapstruct -h

# try to generate sample code

go generate github.com/nikolay-turpitko/x-go/wrapstruct/...

# see original struct with comments and generated wrapper

less $GOPATH/src/github.com/nikolay-turpitko/x-go/wrapstruct/sample/sample.go
less $GOPATH/src/github.com/nikolay-turpitko/x-go/wrapstruct/sample/mystruct_generated.go

```
