package bartender

import (
	"reflect"
)

func (c * Controller) methodInvoker(method string) {	
	argsValue := make([]reflect.Value, len(c.Args))
	
	for i, elem := range c.Args {argsValue[i] = reflect.ValueOf(elem)} 

	// Call action with args
	c.controllerValue.MethodByName(method).Call(argsValue)
}
