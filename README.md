# klog wrapper of zap

inspired by [istio/klog](https://github.com/istio/klog).

## usage

import this package by `github.com/xial-thu/klog`

initialization: for example, in main.go
```golang
    // Unchanged
	klog.InitFlags(nil)
	flag.StringVar(&kubeconfig, "kubeconfig", "", "path to Kubernetes config file")
	flag.Parse()

    // Append a new API
	klog.Singleton()
```

Due to some gaps between klog and zap, parameters shall be converted, and the conversion must be done after `flag.Parse()`. `klog.Singleton()` inits an unique global logger whose configuration is slightly different from default zap production configuration at:

1. `Timekey` is set to "time"
2. `EncodeTime` is set to `ISO8601TimeEncoder`

flags: not all flags defined in klog is supported, or rather say, not all the flags still make sense. Only `alsologtostderr` and `v` is supported currently.

* `v`: still supports `klog.V(2).Info()` syntax. But the level of `klog.Info()` is **INFO**; that `klog.V(3).Info()` is **DEBUG**. If v is set to zero, zap **DEBUG** log will be ignored. Default to **0**
* `alsologtostderr`: default to true

## limitation

1. default context is empty.
2. structured logging still needs modification in business logic.
3. support of flags needs more work.
