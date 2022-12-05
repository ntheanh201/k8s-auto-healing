package handler

import (
	"fmt"
	"log"
	"reflect"
)

type Handler interface {
	Start()
}

type HandlerRegistry struct {
	handlers     map[reflect.Type]Handler
	handlerTypes []reflect.Type
}

func NewHandlerRegistry() *HandlerRegistry {
	return &HandlerRegistry{
		handlers: make(map[reflect.Type]Handler),
	}
}

func (r *HandlerRegistry) RegisterHandler(handler Handler) error {
	kind := reflect.TypeOf(handler)
	if _, ok := r.handlers[kind]; ok {
		return fmt.Errorf("handler already exists: %v", kind)
	}
	r.handlerTypes = append(r.handlerTypes, kind)
	r.handlers[kind] = handler
	return nil
}

func (r *HandlerRegistry) StartAll() {
	log.Printf("Starting %d handlers: %v\n", len(r.handlerTypes), r.handlerTypes)
	for _, kind := range r.handlerTypes {
		log.Printf("Starting service type %v\n", kind)
		go r.handlers[kind].Start()
	}
}

// https://rauljordan.com/2020/03/10/building-a-service-registry-in-go.html
