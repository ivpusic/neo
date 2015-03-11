package neo

import (
	"encoding/json"
	"encoding/xml"
	"net/http"
	"os"
	"path/filepath"
)

// Server response representation.
type Response struct {
	Status  int
	writer  http.ResponseWriter
	request *http.Request
	Cookie  Cookie
	data    []byte
	// defer file sending
	file       string
	_skipFlush bool
}

// making response representation
func makeResponse(request *http.Request, w http.ResponseWriter) *Response {
	response := &Response{
		request: request,
		writer:  w,
		Cookie:  Cookie{},
	}

	return response
}

///////////////////////////////////////////////////////////////////
// Creating Responses
///////////////////////////////////////////////////////////////////

// Will produce JSON string representation of passed object,
// and send it to client
func (r *Response) Json(obj interface{}, status int) {
	res, err := json.Marshal(obj)
	if err != nil {
		log.Error(err.Error())
		return
	}

	r.writer.Header().Set("Content-Type", "application/json")
	r.Raw(res, status)
}

// Will produce XML string representation of passed object,
// and send it to client
func (r *Response) Xml(obj interface{}, status int) {
	res, err := xml.Marshal(obj)
	if err != nil {
		log.Error(err.Error())
		return
	}

	r.writer.Header().Set("Content-Type", "application/xml")
	r.Raw(res, status)
}

// Will send provided Text to client.
func (r *Response) Text(text string, status int) {
	r.Raw([]byte(text), status)
}

// Will look for template, render it, and send rendered HTML to client.
// Second argument is data which will be passed to client.
func (r *Response) Tpl(name string, data interface{}) {
	log.Infof("Rendering template %s", name)

	err := renderTpl(r.writer, name, data)
	if err != nil {
		log.Error(err.Error())
		r.Status = http.StatusNotFound
	} else {
		r.Status = http.StatusOK
	}

	r.skipFlush()
}

// Send Raw data to client.
func (r *Response) Raw(data []byte, status int) {
	r.Status = status
	r.data = data
}

// Write raw response. Implements ResponseWriter.Write.
func (r *Response) Write(b []byte) (int, error) {
	return r.writer.Write(b)
}

// Get Header. Implements ResponseWriter.Header.
func (r *Response) Header() http.Header {
	return r.writer.Header()
}

// Write Header. Implements ResponseWriter.WriterHeader.
func (r *Response) WriteHeader(s int) {
	r.writer.WriteHeader(s)
}

// Checking if file exist.
// todo: consider moving this to utils.go
func (r *Response) fileExists(file string) bool {
	if _, err := os.Stat(file); err == nil {
		return true
	}

	r.Status = http.StatusNotFound
	log.Warnf("cannot find %s file", file)
	return false
}

// Find file, and send it to client.
func (r *Response) File(path string) {
	abspath, err := filepath.Abs(path)

	if err != nil {
		log.Errorf("error while converting %s to absolute path", path)
		return
	}

	if !r.fileExists(abspath) {
		return
	}

	r.file = abspath
}

// Serving static file.
func (r *Response) serveFile(file string) {
	log.Debugf("serving file %s", file)

	r.Status = http.StatusOK
	http.ServeFile(r.writer, r.request, file)
}

// Will be called from ``flush`` Response method if user called ``File`` method.
func (r *Response) sendFile() {
	log.Debugf("sending file %s", r.file)

	base := filepath.Base(r.file)
	r.writer.Header().Set("Content-Disposition", "attachment; filename="+base)
	http.ServeFile(r.writer, r.request, r.file)
}

///////////////////////////////////////////////////////////////////
// Writing Response
///////////////////////////////////////////////////////////////////

// If it is called, Neo will skip calling  ResponseWriter's write method.
// Usecase for this is when we render HTML template for example, because Neo uses Go html/template for
// writing to ResponseWriter.
func (r *Response) skipFlush() {
	r._skipFlush = true
}

// Write result to ResponseWriter.
func (r *Response) flush() {
	if r._skipFlush {
		log.Debug("Already sent. Skipping flushing")
		return
	}

	// set all cookies to response object
	for k, v := range r.Cookie {
		log.Debug(k)
		http.SetCookie(r.writer, v)
	}

	// in case of file call separate function for piping file to client
	if len(r.file) > 0 {
		log.Debugf("found file, sending")
		r.sendFile()
	} else {
		r.writer.WriteHeader(r.Status)
		r.writer.Write(r.data)
	}
}
