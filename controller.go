package bartender

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"reflect"
	"fmt"
	"os"
)

type Controller struct {
	context			*gin.Context
	controllerType		reflect.Type
	controllerValue	reflect.Value
	ControllerName		string
	MethodName		string		
	Action			string		// Full name Ctrlr.Fcn
	TplName 			string
	TplExt 			string
	RenderType		string
	Args				[]interface{}
	HttpStatus		int
	Subdomain			string
}

// add controllers to the bartender instance
func (b *bartender) addController(c *Controller) {
	
	c.controllerType = validateController(c, nil)
	
	c.controllerValue = reflect.New(c.controllerType)
		
	b.controllers[c.Action] = c
	fmt.Println(c.Action, " - ", c.ControllerName, " - ", c.MethodName, " - ", c.HttpStatus)
}

// Create a new controller
func (b *bartender) NewController(action string) *Controller {
	c := new(Controller)
	
	c.context = new(gin.Context)
	
	c.Action = action
	c.actionSplit()
		
	return c
}

// Panics unless validation is correct
func validateController(controller interface{}, parentControllerType reflect.Type) reflect.Type {
	controllerType := reflect.TypeOf(controller)

	if controllerType.Kind() == reflect.Ptr {
		controllerType = controllerType.Elem()
	}

	if controllerType.Kind() != reflect.Struct {
		fmt.Fprintf(os.Stderr,"Controller Type: %s\n", controllerType.Kind())
		panic("Controller needs to be a struct type.")
	}

	if parentControllerType != nil && parentControllerType != controllerType {
		if controllerType.NumField() == 0 {
			panic("Controller needs to have first field be a pointer to parent controller")
		}

		fieldType := controllerType.Field(0).Type

		// Ensure the first field is a pointer to parentControllerType
		if fieldType != reflect.PtrTo(parentControllerType) {
			panic("Controller needs to have first field be a pointer to parent controller")
		}
	}
	
	return controllerType
}

// Fetch the subdomain from the context of the controller
/*
func (c *Controller) getSubdomain(dName string) string {
	host_split := strings.Split(c.Request.Host, "."+dName)
	return host_split[0]
}
*/

func GetSubdomain(dName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		host_split := strings.Split(c.Request.Host, "."+dName)
		subdomain := host_split[0]
		c.Set("subdomain", subdomain)
		
		// Handle the request
		c.Next()
	}
}

func (c *Controller) actionSplit() {
	actionSplit := strings.Split(c.Action, ".")
		
	if len(actionSplit) == 2 {
		c.ControllerName = actionSplit[0]
		c.MethodName = actionSplit[1]
	}
}

// Run the action of the controller
func (c *Controller) Do(method string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// execute action related to controller
		c.HttpStatus = 200
		c.Args = make([]interface{},1)
		c.Args[0] = "Hello World"
			
		c.Render(ctx)
	}
}

func (c *Controller) Render(ctx *gin.Context) {
	switch c.RenderType {
		case "JSON":
			c.context.JSON(c.HttpStatus, c.Args)
		case "XML":
			c.context.XML(c.HttpStatus, c.Args)
		case "HTML":
			c.context.HTML(c.HttpStatus, c.TplName, c.Args)
		case "String":
			c.context.String(c.HttpStatus, c.Args[0].(string), c.Args[1:])
		case "Redirect":
			c.context.Redirect(c.HttpStatus, c.Args[0].(string))
		case "Data":
			// Convert data to byes
			bytes := make([]byte, len(c.Args[1:]))
			for i, elem := range c.Args[1:] {bytes[i] = elem.(byte)} 
			
			c.context.Data(c.HttpStatus, c.Args[0].(string), bytes)

		case "File":
			c.context.File(c.Args[0].(string))

	}
	c.context.String(http.StatusInternalServerError, "Could not render route")
}
