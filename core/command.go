package egos

type Command interface {
	Headers() map[string]interface{}
	SetHeader(string, interface{})
	Command() interface{}
	CommandType() string
}

type CommandDescriptor struct {
	command interface{}
	headers map[string]interface{}
}

func NewCommand(command interface{}) *CommandDescriptor {
	return &CommandDescriptor{
		command: command,
		headers: make(map[string]interface{}),
	}
}

func (c *CommandDescriptor) CommandType() string {
	return typeOf(c.command)
}

func (c *CommandDescriptor) Headers() map[string]interface{} {
	return c.headers
}

func (c *CommandDescriptor) SetHeader(key string, value interface{}) {
	c.headers[key] = value
}

func (c *CommandDescriptor) Command() interface{} {
	return c.command
}
