package goroot

type App struct {
	name     string
	handlers []Handler
}

func (app App) GetName() string {
	return app.name
}

func (app App) RegisterHandler(handler Handler) {
	app.handlers = append(app.handlers, handler)
}

func NewApp(name string) *App {
	return &App{
		name: name,
	}
}
