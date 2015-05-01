package neo

import (
	"os"

	"github.com/ivpusic/golog"
)

var (
	app *Application
	log = golog.GetLogger("github.com/ivpusic/neo")
)

func init() {
	// default log Level
	log.Level = golog.INFO
}

// Getting Neo Application instance. This is singleton function.
// First time when we call this method function will try to parse configuration for Neo application.
// It will look for configuration file provided by ``--config`` CLI argument (if exist).
func App() *Application {
	if app == nil {
		confFile := ""

		for i, arg := range os.Args {
			if arg == "--config" {
				if len(arg) > i+1 {
					confFile = os.Args[i+1]
					break
				}
			}
		}

		app = &Application{}
		app.init(confFile)
	}

	return app
}
