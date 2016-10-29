package sample

// ISample is an interface of a 3rd-party lib, which we have to use and cannot
// cange. It should better be a struct, but another developer decided to
// provide very abstract and decoupled API.
type ISample interface {
	PropA() string
	ValB() int
	SomeC() string
	HTTPServer() string
	Zzzz() string

	// ...many other similar getter-like methods...
}

// myStruct should be an implementation of ISample, but for some reason, it
// should have similarly named fields. For example, it is used to unmarshall
// values from configuration files, or json, so fields should be exported to
// unmarshaller see them.
type myStruct struct {
	PropA string
	// some fields are not exported for some reason, but should be for ISample
	valB  int
	SomeC string
	// other fields are started with abbreviation and could not be simply title-cased
	// or happened to have not exactly the same name, but same value
	httpServerURL string `wrapstruct:"HTTPServer"`
	// yet another happened to have same name, but completely different semantic
	Zzzz interface{} `wrapstruct:"-"`

	// ...many other fields, correspondent to ISample's getters...

	// there are also some fields, used by application and not relevant for ISample
	PropD string `wrapstruct:"-"`
}

// Following line instructs go generate tool to generate myStructWrapper.

//go:generate wrapstruct -src myStruct -dst myStructWrapper -o mystruct_generated.go

// Generated struct does not implement Zzzz(), we need to tweek it.

type myISampleImpl struct {
	myStructWrapper
}

func (myISampleImpl) Zzzz() string {
	return "Magic! ))"
}

// Now use it. A bit awkward, but bearable...

var _ ISample = &myISampleImpl{myStructWrapper{}}
