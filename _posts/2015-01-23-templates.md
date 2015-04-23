---
layout: post
title: Templates
date: 2015-01-22 15:39:41
categories: tutorials
---

Neo uses Go ``html/template`` package as template engine. For detail information about template syntax and feature please read [documentation](http://golang.org/pkg/html/template/).

But let's now focus on what you have to know from Neo's side to use templating.

First what you have to do is to create templates. For each template you have to define name.

Let's define three templates, header footer and index. We will include header and footer templates from index.

```html

{% raw  %}{{define "header"}}{% endraw %}
<p>This is header</p>
{% raw  %}{{end}}{% endraw %}
```

```html
{% raw  %}{{define "footer"}}{% endraw %}
<p>This is footer</p>
{% raw  %}{{end}}{% endraw %}
```

```HTML
{% raw  %}{{define "index"}}{% endraw %}
<!DOCTYPE html>
<html>
    <head>
        <title>Example</title>
    </head>

    <body>
        {% raw  %}{{template "header"}}{% endraw %}
        <div>
            Site content
        </div>
        {% raw  %}{{template "footer"}}{% endraw %}
    </body>
</html>
{% raw  %}{{end}}{% endraw %}
```

Next you have to compile above templates. You have to say to Neo where he can find templates which you will use from you application.

Usually you do that from main function.

```Go
app.Templates(
    "/path/to/templates/*",
    "/another/path/to/templates/template.tpl",
)
```

As you can see you can provide multiple paths where your templates are located.

Then in one of your route handlers you have to call template rendering in order to make HTML and return it to user.

```Go
app.Get("/", func(ctx neo.Ctx) (int, error) {
    return 200, ctx.Res.Tpl("index", nil)
})
```

Second parameter of ``Tpl`` function is data which will be passed into template. So you can make instances of your structs and pass them as second argument. Something like this:

```Go
data := Person{"Some", "Person"}
ctx.Res.Tpl("index", data)
```

And that's it basically about templates.
