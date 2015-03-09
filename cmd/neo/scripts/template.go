package scripts

import (
	"github.com/ivpusic/golog"
	"io"
	"os"
)

type Template interface {
	Build(string)
}

var (
	NEO_ROOT      = os.Getenv("GOPATH") + "/src/github.com/ivpusic/neo"
	NEO_TEMPLATES = NEO_ROOT + "/cmd/neo/templates"
	logger        = golog.GetLogger("github.com/ivpusic/neo")
)

func copyTemplate(template, location string) {
	// copy config template to target directory
	source, err := os.Open(NEO_TEMPLATES + "/" + template)
	checkErr(err)
	defer source.Close()

	dest, err := os.Create(location + "/" + template)
	checkErr(err)
	defer dest.Close()

	_, err = io.Copy(dest, source)
	checkErr(err)
}
