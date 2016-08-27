# extend errors

[referance](http://golangweekly.us1.list-manage.com/track/click?u=0618f6a79d6bb9675f313ceb2&id=f918a227f6&e=0708a4e66b)


### Example errors.Wrap

```golang

err := doSomething()

err = errors.Wrap(err, "in my pkg")

fmt.Println(err)

```
