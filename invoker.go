package bartender

import (
	"reflect"
	"fmt"
)

func (c * Controller) methodInvoker(method string) {
	fmt.Println("sizeof C.Args", len(c.Args))
	argsValue := make([]reflect.Value, len(c.Args))
	
	for i, elem := range c.Args {argsValue[i] = reflect.ValueOf(elem)} 

	fmt.Println(method, "-", argsValue)

	fmt.Println("MethodByName.IsValid() ==", c.controller.MethodByName(method).IsValid())
	// Call action with args
	//c.controllerValue.MethodByName(method).IsValid(argValues)
	//c.controllerValue.MethodByName(method).Call([]reflect.Value{})
	
}
