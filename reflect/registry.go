package reflect

import (
	"reflect"
	"sync"

	"darvaza.org/core"
)

type Registry struct {
	mu sync.RWMutex
}

func (m *Registry) Register(i any) (*Type, error) {
	return m.RegisterType(reflect.TypeOf(i))
}

func (*Registry) RegisterType(rt reflect.Type) (*Type, error) {
	return nil, core.ErrTODO
}

func NewRegistry() (*Registry, error) {
	return nil, core.ErrTODO
}

func MustRegistry() *Registry {
	m, err := NewRegistry()
	if err != nil {
		panic(err)
	}
	return m
}

func Register(t any) (*Type, error) {
	return global.Register(t)
}

var global = MustRegistry()
