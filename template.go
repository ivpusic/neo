package neo

import (
	"errors"
	"html/template"
	"io"
	"path/filepath"
)

// todo: move global Template instance to some local scope
var tpl *template.Template

// Write result of ExecuteTemplate to provided writer. This is usually ResponseWriter.
func renderTpl(writer io.Writer, name string, data interface{}) error {
	if tpl == nil {
		return errors.New("Rendering template error! Please compile your templates.")
	}

	return tpl.ExecuteTemplate(writer, name, data)
}

// Collect all files for come location.
func collectFiles(location string) ([]string, error) {
	log.Debug("collecting templates from " + location)
	result := []string{}

	paths, err := filepath.Glob(location)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	for _, path := range paths {
		isDir, err := isDirectory(path)
		if err != nil {
			log.Error(err.Error())
			return nil, err
		}

		if !isDir {
			result = append(result, path)
		}
	}

	return result, nil
}

// compile templates on provided locations
// method will panic if there is something wrong during template compilation
// this is intended to be called on application startup
func compileTpl(paths ...string) {
	locations := []string{}
	var err error

	for _, path := range paths {
		pths, err := collectFiles(path)
		if err != nil {
			panic(err)
		} else {
			locations = append(locations, pths...)
		}
	}

	tpl, err = template.ParseFiles(locations...)
	if err != nil {
		panic(err)
	}
}

func tplExists(name string) bool {
	return tpl != nil && tpl.Lookup(name) != nil
}
