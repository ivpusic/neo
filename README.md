Neo
====

## Example
```Go
func main() {
    app := neo.App()

    app.Get("/", func(this *neo.Ctx) {
        this.Res.Text("works")
    })

    app.Start(":3000")
}
```
