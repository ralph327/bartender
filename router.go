package bartender  

import (
)

func (b *bartender) AddRoute(method string, route string, action string) {
	// Create a new controller based on the action
	c := b.NewController(action)
	
	// Add the controller to 
	b.addController(c)
	
	
	switch method {
		case "POST":
			b.server.POST(route, b.controllers[c.ControllerName].Do(c.MethodName))
		case "GET":
			b.server.GET(route, b.controllers[c.ControllerName].Do(c.MethodName))
		case "DELETE":
			b.server.DELETE(route, b.controllers[c.ControllerName].Do(c.MethodName))
		case "PATCH":
			b.server.PATCH(route, b.controllers[c.ControllerName].Do(c.MethodName))
		case "PUT":
			b.server.PUT(route, b.controllers[c.ControllerName].Do(c.MethodName))
		case "OPTIONS":
			b.server.OPTIONS(route, b.controllers[c.ControllerName].Do(c.MethodName))
		case "HEAD":
			b.server.HEAD(route, b.controllers[c.ControllerName].Do(c.MethodName))
		case "LINK":
			b.server.LINK(route, b.controllers[c.ControllerName].Do(c.MethodName))
		case "UNLINK":
			b.server.UNLINK(route, b.controllers[c.ControllerName].Do(c.MethodName))
	}
}