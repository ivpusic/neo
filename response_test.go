package neo

import (
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/ivpusic/httpcheck"
	"github.com/stretchr/testify/assert"
)

type testPerson struct {
	Name string
	Age  int
}

func getTestApp(t *testing.T) *httpcheck.Checker {
	app := App()
	app.Templates("./testassets/templates/*")

	app.Get("/json", func(this *Ctx) (int, error) {
		return 200, this.Res.Json(&testPerson{"Some", 30})
	})

	app.Get("/xml", func(this *Ctx) (int, error) {
		return 200, this.Res.Xml(&testPerson{"Some", 30})
	})

	app.Get("/text", func(this *Ctx) (int, error) {
		return 200, this.Res.Text("some text")
	})

	app.Get("/tpl", func(this *Ctx) (int, error) {
		return 200, this.Res.Tpl("testindex", testPerson{"Some", 30})
	})

	app.Get("/cookie", func(this *Ctx) (int, error) {
		this.Res.Cookie.Set("key", "value")
		this.Res.Cookie.Set("key1", "value1")
		return 200, nil
	})

	app.Get("/header", func(this *Ctx) (int, error) {
		this.Res.Header().Set("key", "value")
		this.Res.Header().Set("key1", "value1")
		return 200, nil
	})

	app.Get("/status400", func(this *Ctx) (int, error) {
		return 400, nil
	})

	app.Get("/status200", func(this *Ctx) (int, error) {
		return 200, nil
	})

	app.Get("/file", func(this *Ctx) (int, error) {
		return 200, this.Res.File("./testassets/test.txt")
	})

	app.Get("/fileunknown", func(this *Ctx) (int, error) {
		return 404, this.Res.File("./testassets/test_unkonown_file.txt")
	})

	app.Options("/some/*/:path", func(this *Ctx) (int, error) {
		return 200, this.Res.Text(this.Req.Params.Get("path"))
	})

	app.Options("*", func(this *Ctx) (int, error) {
		return 200, this.Res.Text("allok")
	})

	app.Serve("/assets", "./testassets")

	app.flush()

	return httpcheck.New(t, app)
}

func TestWildcard(t *testing.T) {
	server := getTestApp(t)
	server.Test("OPTIONS", "/some/blabla/value").Check().
		HasStatus(200).
		Cb(func(response *http.Response) {
		body, err := ioutil.ReadAll(response.Body)
		assert.Nil(t, err)
		assert.Equal(t, "value", string(body))
	})

	server.Test("OPTIONS", "/blabla").Check().
		HasStatus(200).
		Cb(func(response *http.Response) {
		body, err := ioutil.ReadAll(response.Body)
		assert.Nil(t, err)
		assert.Equal(t, "allok", string(body))
	})
}

func TestStatus(t *testing.T) {
	server := getTestApp(t)
	server.Test("GET", "/status400").
		Check().
		HasStatus(400)

	server.Test("GET", "/status200").
		Check().
		HasStatus(200)
}

func TestJsonResponse(t *testing.T) {
	server := getTestApp(t)
	server.Test("GET", "/json").
		Check().
		HasHeader("Content-Type", "application/json").
		HasStatus(200).
		HasJson(&testPerson{"Some", 30})
}

func TestXmlResponse(t *testing.T) {
	server := getTestApp(t)
	server.Test("GET", "/xml").
		Check().
		HasHeader("Content-Type", "application/xml").
		HasStatus(200).
		HasXml(&testPerson{"Some", 30})
}

func TestTplResponse(t *testing.T) {
	server := getTestApp(t)
	server.Test("GET", "/tpl").
		Check().
		HasStatus(200).
		Cb(func(response *http.Response) {
		body, err := ioutil.ReadAll(response.Body)
		assert.Nil(t, err)

		html := "\n<html>\n" +
			"<head>\n" +
			"    <title>Test</title>\n" +
			"</head>\n" +
			"<body>\n" +
			"    \n" +
			"Header\n" +
			"\n" +
			"    Body Some\n" +
			"    \n" +
			"Footer\n" +
			"\n" +
			"</body>\n" +
			"</html>\n"

		assert.Equal(t, html, string(body))
	})
}

func TestTextResponse(t *testing.T) {
	server := getTestApp(t)
	server.Test("GET", "/text").
		Check().
		HasHeader("Content-Type", "text/plain; charset=utf-8").
		HasStatus(200).
		Cb(func(response *http.Response) {
		body, err := ioutil.ReadAll(response.Body)
		assert.Nil(t, err)
		assert.Equal(t, "some text", string(body))
	})
}

func TestCookie(t *testing.T) {
	server := getTestApp(t)
	server.Test("GET", "/cookie").
		Check().HasStatus(200).
		HasCookie("key", "value").
		HasCookie("key1", "value1")
}

func TestHeader(t *testing.T) {
	server := getTestApp(t)
	server.Test("GET", "/header").
		Check().HasStatus(200).
		HasHeader("key", "value").
		HasHeader("key1", "value1")
}

func TestServingFile(t *testing.T) {
	server := getTestApp(t)
	server.Test("GET", "/assets/test.txt").Check().
		Cb(func(response *http.Response) {
		body, err := ioutil.ReadAll(response.Body)
		assert.Nil(t, err)
		assert.Equal(t, "file content\n", string(body))
	})

	server.Test("GET", "/assets/test_unknown_file.txt").Check().
		HasStatus(404)
}

func TestSendFile(t *testing.T) {
	server := getTestApp(t)
	server.Test("GET", "/file").Check().
		HasStatus(200).
		Cb(func(response *http.Response) {
		body, err := ioutil.ReadAll(response.Body)
		assert.Nil(t, err)
		assert.Equal(t, "file content\n", string(body))
	})

	server.Test("GET", "/fileunknown").Check().
		HasStatus(404)
}

func TestRouteNotFound(t *testing.T) {
	server := getTestApp(t)
	server.Test("PUT", "/unknown_route").
		Check().
		HasStatus(404)
}
