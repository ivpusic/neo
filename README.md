Neo
====
[![Build Status](https://travis-ci.org/ivpusic/neo.svg?branch=master)](https://travis-ci.org/ivpusic/neo)
[![GoDoc](https://godoc.org/github.com/ivpusic/neo?status.svg)](https://godoc.org/github.com/ivpusic/neo)

Go Web Framework

## Installation

```bash
# framework
go get github.com/ivpusic/neo

# CLI tool
go get github.com/ivpusic/neo/cmd/neo
```

# Documentation
[Project Site](http://ivpusic.github.io/neo)

[API Documentation](http://godoc.org/github.com/ivpusic/neo)

## Example

Create Neo application:
```bash
neo new myapp
cd myapp
```

```Go
package main

import (
    "github.com/ivpusic/neo"
)

func main() {
    app := neo.App()

    app.Get("/", func(ctx *neo.Ctx) (int, error) {
        return 200, ctx.Res.Text("I am Neo Programmer")
    })

    app.Start()
}
```

Run it:
```bash
neo run main.go
```

## License
*MIT*
