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

Useful links on the topic:
https://github.com/golang/example/tree/master/gotypes#struct-types
https://golang.org/pkg/go/types/
https://godoc.org/golang.org/x/tools/go/loader#example-Config-Import
https://github.com/ernesto-jimenez/gogen
