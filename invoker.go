package bartender

import (
	"reflect"
	"fmt"
)

func (c * Controller) methodInvoker(method string) {	
	argsValue := make([]reflect.Value, len(c.Args))
	
	for i, elem := range c.Args {argsValue[i] = reflect.ValueOf(elem)} 

	fmt.Println(method, " - ", argsValue)



	// Call action with args
	//c.controllerValue.MethodByName(method).Call(argsValue)
	//c.controllerValue.MethodByName(method).Call([]reflect.Value{})
	
}
