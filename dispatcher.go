package finkgoes

import (
	"fmt"
)

type Dispatcher interface {
	Dispatch(Command) error
	RegisterHandler(CommandHandler, ...interface{}) error
}

type DispatcherImpl struct {
	handlers map[string]CommandHandler
}

func NewDispatcher() *DispatcherImpl {
	return &DispatcherImpl{
		handlers: make(map[string]CommandHandler),
	}
}

func (d *DispatcherImpl) Dispatch(command Command) error {
	if handler, ok := d.handlers[command.CommandType()]; ok {
		return handler.Handle(command)
	}
	return fmt.Errorf("command handler for commands of type %s does not exist", command.CommandType())
}

func (b *DispatcherImpl) RegisterHandler(handler CommandHandler, commands ...interface{}) error {
	for _, command := range commands {
		typeName := typeOf(command)
		if _, ok := b.handlers[typeName]; ok {
			return fmt.Errorf("duplicate command handler registration for command of type: %s", typeName)
		}
		b.handlers[typeName] = handler
	}
	return nil
}