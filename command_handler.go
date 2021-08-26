package finkgoes

type CommandHandler interface {
	Handle(Command) error
}

type CommandHandlerBase struct {
	next CommandHandler
}