package goroot

type Context struct {
	backpack map[string]interface{}
}

// Get returns the value for the given key
func (ctx Context) Get(key string) interface{} {
	return ctx.backpack[key]
}

// Set sets the value for the given key
func (ctx Context) Set(key string, value interface{}) {
	ctx.backpack[key] = value
}

// NewContext creates a new context
func NewContext() Context {
	return Context{
		backpack: make(map[string]interface{}),
	}
}
