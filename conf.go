package neo

import (
	"time"

	"net/http"

	"github.com/ivpusic/toml"
)

const (
	BYTE     = 1.0
	KILOBYTE = 1024 * BYTE
	MEGABYTE = 1024 * KILOBYTE
	GIGABYTE = 1024 * MEGABYTE
	TERABYTE = 1024 * GIGABYTE
)

///////////////////////////////////////////////////////////////////
// `Hotreload` section
///////////////////////////////////////////////////////////////////

type HotReloadConf struct {
	Suffixes []string
	Ignore   []string
}

type LoggerConf struct {
	Name  string
	Level string
}

///////////////////////////////////////////////////////////////////
// `App` section
///////////////////////////////////////////////////////////////////

type ApplicationConf struct {
	Test         string
	Args         []string
	Addr         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	// default to http.DefaultMaxHeaderBytes
	MaxHeaderBytes int
	// Default 10MB
	MaxBodyBytes int64
	Logger       LoggerConf
}

///////////////////////////////////////////////////////////////////
// `Neo` section
///////////////////////////////////////////////////////////////////

type NeoConf struct {
	Logger LoggerConf
}

///////////////////////////////////////////////////////////////////
// `Global` section
///////////////////////////////////////////////////////////////////

// Neo Application configuration
type Conf struct {
	Hotreload HotReloadConf
	App       ApplicationConf
	Neo       NeoConf
}

func (c *Conf) loadDefaults() {
	// hotreload
	c.Hotreload.Ignore = []string{}
	c.Hotreload.Suffixes = []string{}

	// app
	c.App.Args = []string{}
	c.App.Addr = ":3000"
	c.App.Logger.Level = "INFO"
	c.App.Logger.Name = "application"
	c.App.MaxHeaderBytes = http.DefaultMaxHeaderBytes
	c.App.MaxBodyBytes = 10 * MEGABYTE
	c.App.ReadTimeout = 0
	c.App.WriteTimeout = 0
}

// Will try to parse TOML configuration file.
func (c *Conf) Parse(path string) {
	c.loadDefaults()

	if path == "" {
		log.Info("Loaded configuration defaults")
		return
	}

	if _, err := toml.DecodeFile(path, c); err != nil {
		panic(err)
	}
}
