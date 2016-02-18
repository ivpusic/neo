---
layout: post
title: Static file serving
date: 2015-01-22 15:39:41
categories: tutorials
---

Neo can be fast static file server too.
For example let's take following directory structure.

```
├── assets
│   ├── css
│   │   └── stylesheet.css
│   └── js
│       └── script.js
├── main.go
```

We want to serve ```assets``` directory using Neo. Ok, easy.

```go
app.Serve("/static", "./path/to/assets")
```

Now if we navigate to ```/static/js/script.js``` or ```/static/css/stylesheet.css``` we will get content of those files.
