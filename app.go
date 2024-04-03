package docket

type App struct {
	config Config
}

type Config struct {
}

const Version = "0.1"

func New(config ...Config) *App {
	app := &App{
		config: Config{},
	}

	if len(config) > 0 {
		app.config = config[0]
	}

	return app
}
