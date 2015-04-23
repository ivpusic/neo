package scripts

import (
	"io/ioutil"
	"os"
)

type NeoTemplate struct{}

func (n *NeoTemplate) Build(location string) {
	logger.Info("Preparing Neo application")

	// create target directory
	err := os.MkdirAll(location, 0755)
	checkErr(err)
	err = os.Chdir(location)
	checkErr(err)

	// write template to file
	neoApplicationTemplate := []byte(`
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

	app.Get("/", func(ctx *neo.Ctx) (int, error) {
		return 200, ctx.Res.Text("Works!")
	})

	app.Start()
}`)
	err = ioutil.WriteFile(location+"/main.go", neoApplicationTemplate, 0755)
	checkErr(err)

	copyTemplate("config.toml", location)
	logger.Infof("Done! Enter '%s' and type 'neo run main.go' to run your app.", location)
}
