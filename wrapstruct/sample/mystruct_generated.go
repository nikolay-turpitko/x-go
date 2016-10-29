
/*
* CODE GENERATED AUTOMATICALLY WITH github.com/nikolay-turpitko/x-go/wrapstruct
* THIS FILE SHOULD NOT BE EDITED BY HAND
*/

package sample

type myStructWrapper struct {
	w *myStruct
}

func (w *myStructWrapper) PropA() string {
	return w.w.PropA
}

func (w *myStructWrapper) ValB() int {
	return w.w.valB
}

func (w *myStructWrapper) SomeC() string {
	return w.w.SomeC
}

func (w *myStructWrapper) HTTPServer() string {
	return w.w.httpServerURL
}

