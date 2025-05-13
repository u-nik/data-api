package handlers

import "go.uber.org/zap"

var factories = make(map[string]HandlerFactory)

func RegisterHandler(name string, factory HandlerFactory) {
	if _, exists := factories[name]; exists {
		panic("handler already registered: " + name)
	}
	factories[name] = factory
}

func CreateHandlers(baseHandler BaseHandler) map[string]HandlerInterface {
	handlerMap := make(map[string]HandlerInterface)
	for name, factory := range factories {
		zap.L().Sugar().Infof("Registering handler: %s", name)
		handler := factory(baseHandler)
		handlerMap[name] = handler
	}
	return handlerMap
}
