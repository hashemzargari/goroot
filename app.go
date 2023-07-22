package goroot

// App is the main application.
type App struct {
	name     string
	handlers []Handler
}

// GetName returns the name of the app.
func (app App) GetName() string {
	return app.name
}

// RegisterHandlers registers a handler to the app.
func (app App) RegisterHandlers(handler ...Handler) error {
	for _, h := range handler {
		if !inArray(h, app.handlers) {
			app.handlers = append(app.handlers, h)
		}
	}

	return nil
}

func (app App) Run(grpcAddress string) error {
	// generate protobuf
	// generate grpc
	// generate protoTypes transformers (from proto to handler request/response types)
	// register handlers to grpc server
	// start grpc server
	// TODO
	return nil
}

func inArray(h Handler, handlers []Handler) bool {
	for _, handler := range handlers {
		if getTypeName(h) == getTypeName(handler) {
			return true
		}
	}
	return false
}

// NewApp creates a new app.
func NewApp(name string) *App {
	return &App{
		name: name,
	}
}
