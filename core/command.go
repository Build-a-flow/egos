package egos

type Command interface {
	Command() interface{}
	CommandType() string
}

type CommandDescriptor struct {
	command interface{}
}

func NewCommand(command interface{}) *CommandDescriptor {
	return &CommandDescriptor{
		command: command,
	}
}

func (c *CommandDescriptor) CommandType() string {
	return typeOf(c.command)
}

func (c *CommandDescriptor) Command() interface{} {
	return c.command
}
