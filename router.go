package bartender  

import (
	"strings"
)

func (b *bartender) AddRoute(method string, route string, action string) {
	
	actionSplit := strings.Split(action, ".")
	
	var controllerName, methodName string
	
	if len(actionSplit) == 2 {
		controllerName = actionSplit[0]
		methodName = actionSplit[1]
	}
	
	switch method {
		case "POST":
			b.server.POST(route, b.controllers[controllerName].Do(methodName))
		case "GET":
			b.server.GET(route, b.controllers[controllerName].Do(methodName))
		case "DELETE":
			b.server.DELETE(route, b.controllers[controllerName].Do(methodName))
		case "PATCH":
			b.server.PATCH(route, b.controllers[controllerName].Do(methodName))
		case "PUT":
			b.server.PUT(route, b.controllers[controllerName].Do(methodName))
		case "OPTIONS":
			b.server.OPTIONS(route, b.controllers[controllerName].Do(methodName))
		case "HEAD":
			b.server.HEAD(route, b.controllers[controllerName].Do(methodName))
		case "LINK":
			b.server.LINK(route, b.controllers[controllerName].Do(methodName))
		case "UNLINK":
			b.server.UNLINK(route, b.controllers[controllerName].Do(methodName))
	}
}
