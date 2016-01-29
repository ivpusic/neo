package main

import (
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/ivpusic/golog"
	"github.com/ivpusic/neo"
	"github.com/ivpusic/neo/cmd/neo/scripts"
	"gopkg.in/alecthomas/kingpin.v1"
)

var (
	logger = golog.GetLogger("github.com/ivpusic/neo")
	app    = kingpin.New("Neo", "Command line tool to managing your Neo application")

	// run command
	run        = app.Command("run", "Run Neo application")
	mainFile   = run.Arg("program", "Main file").String()
	configFile = run.Flag("config", "Path to configuration file").Short('c').String()
	verbose    = run.Flag("verbose", "Run in verbose mode").Short('v').Bool()

	// new commmand
	newCmd       = app.Command("new", "Create new Neo application")
	projectName  = newCmd.Arg("name", "Project name").String()
	templateName = newCmd.Flag("tpl", "Create new Neo application based on some template.").Short('t').String()
)

///////////////////////////////////////////////////////////////////
// `run` command
///////////////////////////////////////////////////////////////////

func handleRunCommand() {
	config := &neo.Conf{}

	// check if provided main file exists
	if len(*mainFile) > 0 {
		stat, err := os.Stat(*mainFile)
		if err != nil {
			logger.Errorf("Please provide valid main file. `%s` not found!", *mainFile)
			os.Exit(0)
		} else if stat.IsDir() {
			logger.Errorf("Please provide valid main file. `%s` not valid!", *mainFile)
			os.Exit(0)
		}
	}

	// find project root
	mainFileAbs, err := filepath.Abs(*mainFile)
	if err != nil {
		panic(err)
	}

	// if configuration is not provided by CLI,
	// try to figure out where it could be
	if len(*configFile) == 0 {
		if len(*mainFile) > 0 {
			projectRoot := filepath.Dir(mainFileAbs)
			*configFile = path.Join(projectRoot, "config.toml")
		} else {
			cwd, _ := os.Getwd()
			*configFile = path.Join(cwd, "config.toml")
		}

		// if config file doesn't exists,
		// fallback to empty configuration file -> defaults will be used
		_, err := os.Stat(*configFile)
		if err != nil {
			logger.Warnf("Configuration file at %s not found!", *configFile)
			*configFile = ""
		}
	}

	// parse configuration
	config.Parse(*configFile)

	arguments := []string{}

	// ignore files
	if len(config.Hotreload.Ignore) > 0 {
		arguments = append(arguments, "-i", strings.Join(config.Hotreload.Ignore, ","))
	}

	// suffixes to watch
	if len(config.Hotreload.Suffixes) > 0 {
		arguments = append(arguments, "-s", strings.Join(config.Hotreload.Suffixes, ","))
	}

	// app arguments
	config.App.Args = append(config.App.Args, "--config", *configFile)
	arguments = append(arguments, "-a"+strings.Join(config.App.Args, ","))

	// verbose
	if *verbose {
		arguments = append(arguments, "-v")
	}

	runCmd("rerun", arguments)
}

///////////////////////////////////////////////////////////////////
// `new` command
///////////////////////////////////////////////////////////////////

func handleNewCommand() {
	currentLocation, err := os.Getwd()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(0)
	}
	projectLocation := currentLocation + "/" + *projectName

	var template scripts.Template = nil

	if len(*templateName) == 0 {
		logger.Info("Creating Neo project")
		template = &scripts.NeoTemplate{}
	} else {
		switch *templateName {
		case "angular":
			logger.Info("Creating Neo Angular project")
			template = &scripts.AngularTemplate{}
		case "html":
			logger.Info("Creating Neo HTML project")
			template = &scripts.NeoHtmlTemplate{}
		}
	}

	if template == nil {
		logger.Errorf("Unkonown template %s!", *templateName)
	} else {
		template.Build(projectLocation)
	}
}

///////////////////////////////////////////////////////////////////
// main
///////////////////////////////////////////////////////////////////

func main() {
	logger.Level = golog.INFO

	app.Version("1.0.2")

	// install dependencies
	outputCmd("go", []string{"get", "github.com/ivpusic/rerun"})

	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case run.FullCommand():
		handleRunCommand()
	case newCmd.FullCommand():
		handleNewCommand()
	default:
		app.Usage(os.Stdout)
	}
}
