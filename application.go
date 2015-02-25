package neo

import (
	"net/http"
	"os"

	"github.com/ivpusic/golog"
	"github.com/ivpusic/neo/ebus"
)

// Representing Neo application instance
type Application struct {
	// Event emmiter/receiver
	ebus.EBus
	router
	static *Static

	// Application Configuration parameters
	Conf *Conf

	// Application logger instance
	Logger *golog.Logger
}

// Application initialization.
func (a *Application) init() {
	a.InitEBus()
	a.initRouter()

	// neo logger
	lvl, err := parseLogLevel(a.Conf.Neo.Logger.Level)
	if err != nil {
		log.Warn(err)
	} else {
		log.Level = lvl
	}

	// application logger
	lvl, err = parseLogLevel(a.Conf.App.Logger.Level)
	a.Logger = golog.GetLogger(a.Conf.App.Logger.Name)
	if err != nil {
		log.Warn(err)
	} else {
		a.Logger.Level = lvl
	}
}

// Handler interface ``ServeHTTP`` implementation.
// Method will accept all incomming HTTP requests, and pass requests to appropriate handlers if they are defined.
func (a *Application) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// log all unhandled panic's
	// todo: check performance impact
	defer func() {
		if r := recover(); r != nil {
			a.Emit("error", r)
			panic(r)
		}
	}()

	ctx := makeCtx(req, w)
	request := ctx.Req
	response := ctx.Res

	defer response.flush()

	///////////////////////////////////////////////////////////////////
	// Catch Neo Assertions
	///////////////////////////////////////////////////////////////////
	defer func() {
		if r := recover(); r != nil {
			err, ok := r.(*NeoAssertError)

			if ok {
				response.Raw(err.message, err.status)
				a.Emit("error", r)
			} else {
				// bubble panic
				panic(r)
			}
		}
	}()

	///////////////////////////////////////////////////////////////////
	// Static File Serving
	///////////////////////////////////////////////////////////////////
	if a.static != nil {
		// check if file can be served
		file, err := a.static.match(req.URL.Path)

		if err == nil {
			h := func(ctx *Ctx) {
				response.skipFlush()
				response.serveFile(file)
			}

			fn := compose(merge(a.middlewares, []appliable{handler(h)}))
			fn(ctx)
			return
		}

		log.Debug("result not found in static")
	}

	///////////////////////////////////////////////////////////////////
	// Route Matching
	///////////////////////////////////////////////////////////////////
	route, err := a.match(request)

	if err != nil {
		log.Debugf("route %s not found", req.URL.Path)

		// dummy route handler
		h := func(ctx *Ctx) {
			response.Status = http.StatusNotFound
		}

		compose(merge(a.middlewares, []appliable{handler(h)}))(ctx)
	} else {
		route.fnChain(ctx)
	}
}

// Starting application instance. This will run application on port defined by configuration.
func (a *Application) Start() {
	// load the configuration file provided by --config flag if the user hasn't
	// set the config file via app.SetConfigFile() method.
	if a.Conf == nil {
		confFile := ""

		for i, arg := range os.Args {
			if arg == "--config" {
				if len(arg) > i+1 {
					confFile = os.Args[i+1]
					break
				}
			}
		}
		a.SetConfigFile(confFile)
	}

	// initialize the app
	a.init()

	a.flush()

	log.Infof("Starting application on address `%s`", a.Conf.App.Addr)
	err := http.ListenAndServe(a.Conf.App.Addr, a)
	if err != nil {
		panic(err.Error())
	}
}

// Defining paths for serving static files. For example if we say:
// a.Serve("/some", "./mypath")
// then if we require ``/some/js/file.js``, Neo will look for file at ``./mypath/js/file.js``.
func (a *Application) Serve(url string, path string) {
	if a.static == nil {
		log.Debug("creating `Static` instance")

		a.static = &Static{}
		a.static.init()
	}

	a.static.Serve(url, path)
}

// If you are planning to return templates from Neo route handler, then you have to compile them.
// This method will accept list of paths/files and compile them.
// You can use also paths with wildcards (example: /some/path/*).
func (a *Application) Templates(templates ...string) {
	compileTpl(templates...)
}

// Making new region instance. You can create multiple regions.
func (a *Application) Region() *Region {
	return a.makeRegion()
}

///////////////////////////////////////////////////////////////////
// Configuration
///////////////////////////////////////////////////////////////////

// SetConfigFile lets you optionally set custom config file path.
func (a *Application) SetConfigFile(confFile string) {
	if a.Conf == nil {
		a.Conf = new(Conf)
		a.Conf.Parse(confFile)
	}
}
