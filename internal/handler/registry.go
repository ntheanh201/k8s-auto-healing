package handler

import (
	"fmt"
	"log"
	"reflect"
)

type Handler interface {
	StartNewJob()
}

type Registry struct {
	handlers     map[reflect.Type]Handler
	handlerTypes []reflect.Type
}

func NewHandlerRegistry() *Registry {
	return &Registry{
		handlers: make(map[reflect.Type]Handler),
	}
}

func (r *Registry) RegisterHandler(handler Handler) error {
	kind := reflect.TypeOf(handler)
	if _, ok := r.handlers[kind]; ok {
		return fmt.Errorf("handler already exists: %v", kind)
	}
	r.handlerTypes = append(r.handlerTypes, kind)
	r.handlers[kind] = handler
	return nil
}

func (r *Registry) StartAll() {
	log.Printf("Starting %d handlers: %v\n", len(r.handlerTypes), r.handlerTypes)
	for _, kind := range r.handlerTypes {
		log.Printf("Starting handler type %v\n", kind)
		go r.handlers[kind].StartNewJob()
	}
}

// https://rauljordan.com/2020/03/10/building-a-service-registry-in-go.html
