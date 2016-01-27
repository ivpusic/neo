---
layout: post
title: Configuration
date: 2015-01-22 15:39:41
categories: tutorials
---

## With toml

You can configure your application using ``toml`` based configuration file. By default Neo will look for configuration file with name ``config.toml`` in the same directory as application ``main.go`` file.

You can override this behaviour by passing configuration file path using ``-c`` or ``--config`` CLI option.

```bash
neo run --config /some/custom/location/conf.toml main.go
```

Here is full configuration file. All options are optional.
Values bellow are default values. If for example ``addr`` options is missing in configuration file, Neo will set it automatically to it's default value, in this case ``:3000``.

```bash
#
# settings related to recompiling and reruning app when source changes
#
[hotreload]
  # file suffixes to watch for changes
  suffixes = [".go"]

  # files/directories to ignore
  ignore = []

#
# general application settings
#
[app]
  # additional application arguments
  args = []
  addr = ":3000"

  [app.logger]
    level = "DEBUG"
    name = "application"

#
# neo settings
#
[neo]
  [neo.logger]
  level = "INFO"
```

This file will be automatically generated when you create Neo application using ``new`` command.

## With Go

If you don't like configuration files, you can do it from Go code.

Example (listen on 9000):

```Golang
func main() {
    myapp := neo.App()

    myapp.Get("/", func(ctx *neo.Ctx) (int, error) {
        return 200, ctx.Res.Text("I am Neo Programmer")
    })
    
    myapp.Conf.App.Addr = ":9000"
    myapp.Start()
}
```

For more information about `Conf` structure check out api docs here. https://godoc.org/github.com/ivpusic/neo#Conf 
