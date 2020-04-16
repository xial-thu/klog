# klog wrapper of zap

inspired by [istio/klog](https://github.com/istio/klog).

## usage

### import

import this package by `github.com/xial-thu/klog`

### initialization

for example, in main.go

```golang
// Unchanged
klog.InitFlags(nil)
flag.StringVar(&kubeconfig, "kubeconfig", "", "path to Kubernetes config file")
flag.Parse()

// Append a new API
klog.Singleton()
```

If `Singleton()` is not called, the default global no-ops logger will work, which means you are not able to see any real log.

Due to some gaps between klog and zap, parameters shall be converted, and the conversion must be done after `flag.Parse()`. `klog.Singleton()` inits an unique global logger whose configuration is slightly different from default zap production configuration at:

1. `Timekey` is set to "time"
2. `EncodeTime` is set to `ISO8601TimeEncoder`

### flags

Not all flags defined in klog is supported, or rather say, not all the flags still make sense. Only `alsologtostderr` and `v` is supported currently.

* `v`: still supports `klog.V(2).Info()` syntax. But the level of `klog.Info()` is **INFO**; that `klog.V(3).Info()` is **DEBUG**. If v is set to zero, zap **DEBUG** log will be ignored. Default to **0**
* `alsologtostderr`: default to true. If set to false, INFO and DEBUG log will only output to stdout

### structured logging

There're 3 APIs:

* `With()`: parse each field and value from input. `WithFields(struct{A string}{"hi"})` will output `"A":"hi"`. If you care the fields in your struct and hope to extract them, use `With()`
* `WithAll()`: sugar of `zap.Any()`. e.g. `WithFields(struct{A string}{"hi"})` will output `"":{"A":"hi"}`. If you want to record the name of your struct, use `WithAll()`
* `WithFields()`: e.g. `WithFields("ID", 1, "name": "hi")`, just another sugar of `sugar.With()`

Tips of `With()`:

1. Only struct or map will be accepted
2. If arg is a struct, only exported field(which means "FieldName", not "fieldname") will be logged
3. If arg is a map, only accept maps whose key is string
4. If you have nested structs or maps, consider spliting them into several anomony parts
5. If you don't want some fields to be logged automatically, unexport them

Some examples of `With()`:

```golang
type S struct {
	A int
	B string
}
type Q struct {
	D ID
}
s := S{
	A: 10,
	B: "abc",
}
q := Q{
	D: ID(1),
}
c := ""
    
// struct args
With(s).Info(c) // "A":10,"B":"abc"
With(s, q).Info(c) // "A":10,"B":"abc","D":1
// anomony
With(struct {
	A int
	B int
}{1, 2}).Info(c) // "A":1,"B":2
// map args
With(map[string]int{"A":1}).Info(c) // "A":1
```

### performance of `With()`

|  case | ns/op  | B/op | allocs/op |
|  ---- | ----  | ---- | ---- |
| With every time | 3121 |2884 | 20|
| WithField every time | 2239 |1568|9|
|WithAll|3215|2642|14|
|With once|573|7|1|

Due to reflect, `With()` is slow. It indicates us that it's better to write like this instead of parsing interfaces every time.

```golang
newLogger := klog.With(something)
newLogger.Info(something)
```

## limitation

1. default field is empty.
2. structured logging still needs modification in business logic.
3. support of flags needs more work.
