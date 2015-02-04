package neo

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
)

type Static struct {
	serving map[string]string
}

func (s *Static) init() {
	s.serving = map[string]string{}
}

// Trying to find file from passed url.
// If you have called ``app.Assert("/some", "./dir")``, then if url is ``/some/path/file.txt``,
// then method will look for file at ``./dir/path/file.txt``.
func (s *Static) match(url string) (string, error) {
	log.Debug("trying to match static url: " + url)

	if url == "/" {
		url = "/index.html"
	}

	for k, v := range s.serving {
		if strings.Index(url, k) == 0 {
			log.Debugf("replacing %s with %s", url, k)
			file := strings.Replace(url, k, "", 1)

			if k == "/" {
				file = k + file
			}

			path := v + file
			log.Debug("found possible result. Path: " + path)

			if _, err := os.Stat(path); err == nil {
				// if !stat.IsDir() {
				return path, nil
				// }
			}
		}
	}

	return "", errors.New("cannot match " + url + " in static")
}

// Serving static files from some directory.
// If provided url is for example ``/assets``,
// then all requests starting with ``/assets`` will be served from directory
// provided by ``path`` parameter.
func (s *Static) Serve(url string, path string) {
	path, err := filepath.Abs(path)
	if err != nil {
		panic("Cannot make absolute path for " + path)
	}

	log.Debug("Serving static at " + path)

	s.serving[url] = path
}
