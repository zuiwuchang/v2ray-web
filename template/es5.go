package template

import (
	"fmt"

	"github.com/dop251/goja"
)

func renderES5(text string, data interface{}) (s string, e error) {
	vm := goja.New()
	_, e = vm.RunString(text)
	if e != nil {
		return
	}
	render, ok := goja.AssertFunction(vm.Get(`render`))
	if !ok {
		e = fmt.Errorf(`not found function render`)
		return
	}
	result, e := render(goja.Undefined(), vm.ToValue(data))
	if e != nil {
		return
	}
	json := vm.GlobalObject().Get("JSON").ToObject(vm)
	stringify, _ := goja.AssertFunction(json.Get(`stringify`))
	v, _ := stringify(goja.Undefined(), result, goja.Undefined(), vm.ToValue("\t"))
	s = v.String()
	return
}
