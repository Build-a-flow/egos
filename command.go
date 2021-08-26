package finkgoes

type Command interface {
	AggregateID() string
	Headers() map[string]interface{}
	SetHeader(string, interface{})
	Command() interface{}
	CommandType() string
}

type CommandDescriptor struct {
	id      string
	command interface{}
	headers map[string]interface{}
}

func NewCommand(aggregateID string, command interface{}) *CommandDescriptor {
	return &CommandDescriptor{
		id:      aggregateID,
		command: command,
		headers: make(map[string]interface{}),
	}
}

func (c *CommandDescriptor) CommandType() string {
	return typeOf(c.command)
}

func (c *CommandDescriptor) AggregateID() string {
	return c.id
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