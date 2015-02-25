package neo

import (
	"github.com/ivpusic/toml"
)

///////////////////////////////////////////////////////////////////
// `Hotreload` section
///////////////////////////////////////////////////////////////////

type HotReload struct {
	Command string
	Watch   []string
	Ignore  []string
}

type LoggerSettings struct {
	Name  string
	Level string
}

type AppModeSettings struct {
	Addr string
}

///////////////////////////////////////////////////////////////////
// `App` section
///////////////////////////////////////////////////////////////////

type ApplicationSettings struct {
	Addr   string
	Logger LoggerSettings
}

///////////////////////////////////////////////////////////////////
// `Neo` section
///////////////////////////////////////////////////////////////////

type NeoSettings struct {
	Logger LoggerSettings
}

///////////////////////////////////////////////////////////////////
// `Global` section
///////////////////////////////////////////////////////////////////

type Conf struct {
	Hotreload HotReload
	App       ApplicationSettings
	Neo       NeoSettings
}

func (c *Conf) loadDefaults() {
	// hotreload
	c.Hotreload.Ignore = []string{}
	c.Hotreload.Watch = []string{"."}

	// app
	c.App.Addr = ":3000"
	c.App.Logger.Level = "INFO"
	c.App.Logger.Name = "application"
}

// Will try to parse TOML configuration file.
func (c *Conf) Parse(path string) {
	if path == "" {
		log.Warn("Configuration file not provided! Using defaults.")
		c.loadDefaults()
		return
	}

	if _, err := toml.DecodeFile(path, c); err != nil {
		log.Error("An error occured while parsing config file, panicking!")
		panic(err)
	}
}
