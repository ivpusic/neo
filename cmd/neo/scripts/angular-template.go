package scripts

import (
	"io/ioutil"
	"os"
	"os/exec"
)

func checkErr(err error) {
	if err != nil {
		logger.Panic(err)
	}
}

func pipe(cmd *exec.Cmd) *exec.Cmd {
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd
}

type AngularTemplate struct{}

func (a *AngularTemplate) Build(location string) {
	os.MkdirAll(location, 0755)
	client := location + "/client"
	server := location + "/server"

	// clone angular template
	logger.Info("Preparing angular template")
	err := pipe(exec.Command("git", "clone", "https://github.com/angular/angular-seed.git", location+"/client")).Run()
	checkErr(err)

	// install angular deps
	err = os.Chdir(client)
	checkErr(err)
	logger.Info("Installing Angular application dependencies")
	err = pipe(exec.Command("npm", "install")).Run()

	// create server directory
	os.Chdir(server)
	err = os.Mkdir(server, 0755)
	checkErr(err)
	os.Chdir(server)

	// create Neo template
	logger.Info("Preparing Neo application")
	neoAppTemplate := []byte(`
package main

import (
	"github.com/ivpusic/golog"
	"github.com/ivpusic/neo"
	"github.com/ivpusic/neo/middlewares/logger"
)

var (
	log = golog.GetLogger("application")
)

func main() {
	log.Info("Regards from Neo")

	app := neo.App()
	app.Use(logger.Log)

	app.Serve("/", "../client/app/")
	app.Start()
}`)
	err = ioutil.WriteFile(server+"/main.go", neoAppTemplate, 0755)
	checkErr(err)

	// copy config
	copyTemplate("config.toml", server)

	logger.Infof("Done! Enter '%s' and type 'neo run main.go' to run your app.", server)
}
