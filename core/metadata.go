package egos

import "context"

type Metadata struct {
	data map[string]interface{}
}

func NewMetadata(data *map[string]interface{}) *Metadata {
	return &Metadata{
		data: *data,
	}
}

func (m *Metadata) Add(key string, value interface{}) {
	m.data[key] = value
}

func (m Metadata) Get(key string) interface{} {
	return m.data[key]
}

func (m Metadata) All() map[string]interface{} {
	return m.data
}

func (m *Metadata) Context(ctx context.Context) context.Context {
	return context.WithValue(ctx, "metadata", m)
}
