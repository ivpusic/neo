---
layout: post
title: Introduction
date: 2015-01-22 15:39:41
categories: tutorials
---

Here we are assuming that you already have installed Neo and you have working Go environment.

This will be step by step tutorial. At the end we will show complete example.
So first declare your package and import Neo.

```Go
package main

import "github.com/ivpusic/neo"
```

After that declare ``main`` function from where we will start our Web Server.

```Go
func main() {
    // Neo code will go here
}
```

We have basic structure for now. Let's add some Neo code.
You have to create Neo application instance using:

```Go
app := neo.App()
```
You will use ``app`` variable to access many things from Neo. One of most important things in Neo is declaring routes. You define new route by calling corresponding function and by passing path and route handler implementation.

```Go
app.Get("/", func(this *neo.Ctx) {
    this.Res.Text("I am Neo Programmer", 200)
})
```
If you want to declare ``POST`` route to for example ``/some/route`` route, you can use something like this:

```Go
app.Post("/some/route", func(this *neo.Ctx) {
    // your route handler implementation
})
```

At the end of this simple introduction you have to start server on some address and port. We didn't provide any configuration yet, and default address is ``localhost:3000``.

```Go
app.Start()
```

Complete example would be:

```Go
package main

import (
    "github.com/ivpusic/neo"
)

func main() {
    app := neo.App()

    app.Get("/", func(this *neo.Ctx) {
        this.Res.Text("I am Neo Programmer", 200)
    })

    app.Start()
}
```

Save this into ``main.go`` for example.

Recommended way of running Neo application is using Neo CLI tool.

```bash
neo run main.go
```

For more information about Neo CLI tool please look at <a href="{{site.url}}/tutorials/2015/01/22/clitool.html">CLI Tool documentation</a>.

Simple, isn't it? :)