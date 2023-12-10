package util

import (
	"sync"
)

// Properties represents a simple key-value store
type Properties struct {
	mu         sync.RWMutex
	properties map[string]string
}

// NewProperties creates a new Properties instance
func NewProperties() *Properties {
	return &Properties{
		properties: make(map[string]string),
	}
}

// SetProperty sets a property with the given key and value
func (p *Properties) SetProperty(key, value string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.properties[key] = value
}

// GetProperty retrieves the value of a property with the given key
func (p *Properties) GetProperty(key string) (string, bool) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	value, ok := p.properties[key]
	return value, ok
}
