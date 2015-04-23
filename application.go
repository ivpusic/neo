package neo

import (
	"fmt"
	"github.com/ivpusic/golog"
	"github.com/ivpusic/neo/ebus"
	"net/http"
	"runtime/debug"
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
func (a *Application) init(confFile string) {
	a.InitEBus()
	a.initRouter()
	a.initConf(confFile)

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
			fmt.Printf("%v", r)
			a.Emit("error", r)
			log.Panic(r)
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

			// if there is started transaction, do rollback
			if ctx.Tx != nil {
				log.Error("Will rollback transaction")
				txErr := ctx.Tx.Rollback().Error
				if txErr != nil {
					log.Errorf("Error while transaction rollback! %s", txErr.Error())
				}
			}

			if ok {
				response.Raw(err.message)
				response.Status = err.status
			} else {
				log.Errorf("%v", r)
				debug.PrintStack()
				response.Status = 500
			}

			a.Emit("error", r)
		}
	}()

	///////////////////////////////////////////////////////////////////
	// Static File Serving
	///////////////////////////////////////////////////////////////////
	if a.static != nil {
		// check if file can be served
		file, err := a.static.match(req.URL.Path)

		if err == nil {
			h := func(ctx *Ctx) (int, error) {
				response.skipFlush()
				return 200, response.serveFile(file)
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
		log.Debugf("route %s not found. Error: %s", req.URL.Path, err.Error())

		// dummy route handler
		h := func(ctx *Ctx) (int, error) {
			return http.StatusNotFound, nil
		}

		compose(merge(a.middlewares, []appliable{handler(h)}))(ctx)
	} else {
		route.fnChain(ctx)
	}
}

// Starting application instance. This will run application on port defined by configuration.
func (a *Application) Start() {
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

// Parsing configuration file for this application instance.
// If configuration file is not provided, then application will use default settings.
func (a *Application) initConf(confFile string) {
	if a.Conf == nil {
		a.Conf = new(Conf)
		a.Conf.Parse(confFile)
	}
}
