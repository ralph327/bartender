package bartender

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"reflect"
)

type Controller struct {
	*gin.Context
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
}

// Create a new controller
func (b *bartender) NewController(action string) *Controller {
	c := new(Controller)
	
	c.Action = action
	c.actionSplit()
		
	return c
}

// Panics unless validation is correct
func validateController(controller interface{}, parentControllerType reflect.Type) reflect.Type {
	controllerType := reflect.TypeOf(controller)

	if controllerType.Kind() != reflect.Struct {
		panic("Controller needs to be a struct type")
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
	// execute action related to controller
	
	c.methodInvoker(method)
	
	return c.Render()
}

func (c *Controller) Render() gin.HandlerFunc {
	switch c.RenderType {
		case "JSON":
			return func(*gin.Context){
				c.JSON(c.HttpStatus, c.Args)
			}
		case "XML":
			return func(*gin.Context){
				c.XML(c.HttpStatus, c.Args)
			}
		case "HTML":
			return func(*gin.Context){
				c.HTML(c.HttpStatus, c.TplName, c.Args)
			}
		case "String":
			return func(*gin.Context){
				c.String(c.HttpStatus, c.Args[0].(string), c.Args[1:])
			}
		case "HTMLString":
			return func(*gin.Context){
				c.HTMLString(c.HttpStatus, c.Args[0].(string), c.Args[1:])
			}
		case "Redirect":
			return func(*gin.Context){
				c.Redirect(c.HttpStatus, c.Args[0].(string))
			}
		case "Data":
			// Convert data to byes
			bytes := make([]byte, len(c.Args[1:]))
			for i, elem := range c.Args[1:] {bytes[i] = elem.(byte)} 
			
			return func(*gin.Context){
				c.Data(c.HttpStatus, c.Args[0].(string), bytes)
			}
		case "File":
			return func(*gin.Context){
				c.File(c.Args[0].(string))
			}
	}
	return func(*gin.Context){
		c.String(http.StatusInternalServerError, "Could not render route")
	}
}
