---
layout: post
title: Middlewares
date: 2015-01-22 15:39:41
categories: tutorials
---

## Basics

Neo has built in powerful concept of middlewares. This enables you to define set of actions which can be executed before and after set of some routes or even single route.

So how it works?
It can work on multiple levels. You can define middlewares which will be executed before all routes, before some set of routes or before single route.

Let's see how we can define middleware which is called before all routes.

```Go
app.Use(func(this *neo.Ctx, next neo.Next) {
    // middleware implementation
})
```
This piece of code will be called every time when you call something on your Web Application.

Middleware can be every function with **func(*neo.Ctx, neo.Next)** signature. First parameter is the same as in route, and the second one is interesting. It is actually function without any parameters, and if you call it, you will basically invoke rest of your application stack. If you don't call it, then Neo will assume that request lifetime should end, and it will return result to client.

All code after your ``next`` function call will be invoked once after all middlewares bellow this one and route handler are invoked.

## What are use cases for this?
### Global scope

Use case can be checking if user is logged in for example. If user is logged in you will invoke ``next`` function, if it is not, you can set status to 401 without calling ``next`` function.

Also good usecase can be logging all traffic on your Web Application.

```Go
app.Use(func (this *neo.Ctx, next neo.Next) {
    start = time.Now()
    fmt.Printf("--> [Req] %s to %s", this.Req.Method, this.Req.URL.Path)

    next()

    elapsed := float64(time.Now().Sub(start) / time.Millisecond)
    fmt.Printf("<-- [Res] (%d) %s to %s Took %vms", this.Res.Status, this.Req.Method, this.Req.URL.Path, elapsed)
})
```

### Route scope
As we mentioned before we can define middleware for single route. Good usecase for this can be authentication, authorization, etc.

Simple example is here:

```Go
route := app.Get("/some/protected/route", func(this *neo.Ctx) {
    // implementation
})

route.Use(func(this *neo.Ctx, next neo.Next) {
    if authorized {
        next()
    } else {
        this.Res.Status = 401
    }
})
```