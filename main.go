package finkgoes

func main()  {

	commandHandler := &TestCommandHandlers{}
	dispatcher := NewDispatcher()
	_ = dispatcher.RegisterHandler(commandHandler, TestCommand{})
}

type TestCommandHandlers struct {
}

func (h *TestCommandHandlers) Handle(command Command) error {
	return nil
}


type TestCommand struct {
}
