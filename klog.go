/*
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package klog

import (
	"flag"
	"fmt"
	"os"
	"reflect"
	"sync"
	"sync/atomic"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Level is a shim
type Level int32

// Verbose is a shim
type Verbose bool

type config struct {
	zapConfig       zap.Config
	v               int
	alsologtostderr bool
}

// Klogger wraps a sugarlogger
type Klogger struct {
	sugar *zap.SugaredLogger
	level Level
}

var (
	klogger *Klogger
	c       config
	once    sync.Once
)

// init as the global no-ops logger so that unit test will not crash
func init() {
	klogger = &Klogger{
		sugar: zap.S(),
		level: 0,
	}
}

// Singleton inits an unique logger
func Singleton() *Klogger {
	once.Do(func() {
		c.zapConfig = zap.NewProductionConfig()

		// change time from ns to formatted
		c.zapConfig.EncoderConfig.TimeKey = "time"
		c.zapConfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

		// due to gaps between zap and klog
		if c.v > 0 {
			c.zapConfig.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
			klogger.level.Set(Level(c.v))
		}
		if !c.alsologtostderr {
			c.zapConfig.OutputPaths = []string{"stdout"}
		}

		// trace the real source caller due to munual inline is not supported
		zlogger, err := c.zapConfig.Build(zap.AddCallerSkip(1))
		if err != nil {
			panic(err)
		}
		klogger.sugar = zlogger.Sugar()
		Infof("init zap logger...")
	})
	return klogger
}

// InitFlags is a shim
func InitFlags(flagset *flag.FlagSet) {
	flag.IntVar(&c.v, "v", 0, "verbosity of info log")
	flag.BoolVar(&c.alsologtostderr, "alsologtostderr", true, "also write logs to stderr, default to true")
}

// Flush is a shim
func Flush() {
	klogger.sugar.Sync()
}

// Set sets the value of the Level.
func (l *Level) Set(val Level) {
	atomic.StoreInt32((*int32)(l), int32(val))
}

// get returns the value of the Level.
func (l *Level) get() Level {
	return Level(atomic.LoadInt32((*int32)(l)))
}

// V is a shim
func V(level Level) Verbose {
	return Verbose(level <= klogger.level.get())
}

// V is a shim
func (k *Klogger) V(level Level) Verbose {
	return Verbose(level <= k.level.get())
}

// Info is a shim
//go:noinline
func (v Verbose) Info(args ...interface{}) {
	if v {
		klogger.sugar.Debug(args...)
	}
}

// Infoln is a shim
//go:noinline
func (v Verbose) Infoln(args ...interface{}) {
	if v {
		s := fmt.Sprint(args...)
		klogger.sugar.Debug(s, "\n")
	}
}

// Infof is a shim
//go:noinline
func (v Verbose) Infof(format string, args ...interface{}) {
	if v {
		klogger.sugar.Debugf(format, args...)
	}
}

// Info is a shim
//go:noinline
func Info(args ...interface{}) {
	klogger.sugar.Info(args)
}

// Info is a shim
//go:noinline
func (k *Klogger) Info(args ...interface{}) {
	k.sugar.Info(args...)
}

// InfoDepth is a shim
//go:noinline
func InfoDepth(depth int, args ...interface{}) {
	klogger.sugar.Info(args...)
}

// InfoDepth is a shim
//go:noinline
func (k *Klogger) InfoDepth(depth int, args ...interface{}) {
	k.sugar.Info(args...)
}

// Infoln is a shim
//go:noinline
func Infoln(args ...interface{}) {
	s := fmt.Sprint(args...)
	klogger.sugar.Info(s, "\n")
}

// Infoln is a shim
//go:noinline
func (k *Klogger) Infoln(args ...interface{}) {
	s := fmt.Sprint(args...)
	k.sugar.Info(s, "\n")
}

// Infof is a shim
//go:noinline
func Infof(format string, args ...interface{}) {
	klogger.sugar.Infof(format, args...)
}

// Infof is a shim
//go:noinline
func (k *Klogger) Infof(format string, args ...interface{}) {
	k.sugar.Infof(format, args...)
}

// Warning is a shim
//go:noinline
func Warning(args ...interface{}) {
	klogger.sugar.Warn(args...)
}

// Warning is a shim
//go:noinline
func (k *Klogger) Warning(args ...interface{}) {
	k.sugar.Warn(args...)
}

// WarningDepth is a shim
//go:noinline
func WarningDepth(depth int, args ...interface{}) {
	klogger.sugar.Warn(args...)
}

// WarningDepth is a shim
//go:noinline
func (k *Klogger) WarningDepth(depth int, args ...interface{}) {
	k.sugar.Warn(args...)
}

// Warningln is a shim
//go:noinline
func Warningln(args ...interface{}) {
	s := fmt.Sprint(args...)
	klogger.sugar.Warn(s, "\n")
}

// Warningln is a shim
//go:noinline
func (k *Klogger) Warningln(args ...interface{}) {
	s := fmt.Sprint(args...)
	k.sugar.Warn(s, "\n")
}

// Warningf is a shim
//go:noinline
func Warningf(format string, args ...interface{}) {
	klogger.sugar.Warnf(format, args...)
}

// Warningf is a shim
//go:noinline
func (k *Klogger) Warningf(format string, args ...interface{}) {
	k.sugar.Warnf(format, args...)
}

// Error is a shim
//go:noinline
func Error(args ...interface{}) {
	klogger.sugar.Error(args...)
}

// Error is a shim
//go:noinline
func (k *Klogger) Error(args ...interface{}) {
	k.sugar.Error(args...)
}

// ErrorDepth is a shim
//go:noinline
func ErrorDepth(depth int, args ...interface{}) {
	klogger.sugar.Error(args...)
}

// ErrorDepth is a shim
//go:noinline
func (k *Klogger) ErrorDepth(depth int, args ...interface{}) {
	k.sugar.Error(args...)
}

// Errorln is a shim
//go:noinline
func Errorln(args ...interface{}) {
	s := fmt.Sprint(args...)
	klogger.sugar.Error(s, "\n")
}

// Errorln is a shim
//go:noinline
func (k *Klogger) Errorln(args ...interface{}) {
	s := fmt.Sprint(args...)
	k.sugar.Error(s, "\n")
}

// Errorf is a shim
//go:noinline
func Errorf(format string, args ...interface{}) {
	klogger.sugar.Errorf(format, args...)
}

// Errorf is a shim
//go:noinline
func (k *Klogger) Errorf(format string, args ...interface{}) {
	k.sugar.Errorf(format, args...)
}

// Fatal is a shim
//go:noinline
func Fatal(args ...interface{}) {
	klogger.sugar.Error(args...)
	os.Exit(255)
}

// Fatal is a shim
//go:noinline
func (k *Klogger) Fatal(args ...interface{}) {
	k.sugar.Error(args...)
	os.Exit(255)
}

// FatalDepth is a shim
//go:noinline
func FatalDepth(depth int, args ...interface{}) {
	klogger.sugar.Error(args...)
	os.Exit(255)
}

// FatalDepth is a shim
//go:noinline
func (k *Klogger) FatalDepth(depth int, args ...interface{}) {
	k.sugar.Error(args...)
	os.Exit(255)
}

// Fatalln is a shim
//go:noinline
func Fatalln(args ...interface{}) {
	s := fmt.Sprint(args...)
	klogger.sugar.Error(s, "\n")
	os.Exit(255)
}

// Fatalln is a shim
//go:noinline
func (k *Klogger) Fatalln(args ...interface{}) {
	s := fmt.Sprint(args...)
	k.sugar.Error(s, "\n")
	os.Exit(255)
}

// Fatalf is a shim
//go:noinline
func Fatalf(format string, args ...interface{}) {
	klogger.sugar.Errorf(format, args...)
	os.Exit(255)
}

// Fatalf is a shim
//go:noinline
func (k *Klogger) Fatalf(format string, args ...interface{}) {
	k.sugar.Errorf(format, args...)
	os.Exit(255)
}

// Exit is a shim
//go:noinline
func Exit(args ...interface{}) {
	klogger.sugar.Error(args...)
	os.Exit(1)
}

// Exit is a shim
//go:noinline
func (k *Klogger) Exit(args ...interface{}) {
	k.sugar.Error(args...)
	os.Exit(1)
}

// ExitDepth is a shim
//go:noinline
func ExitDepth(depth int, args ...interface{}) {
	klogger.sugar.Error(args...)
	os.Exit(1)
}

// ExitDepth is a shim
//go:noinline
func (k *Klogger) ExitDepth(depth int, args ...interface{}) {
	k.sugar.Error(args...)
	os.Exit(1)
}

// Exitln is a shim
//go:noinline
func Exitln(args ...interface{}) {
	s := fmt.Sprint(args...)
	klogger.sugar.Error(s, "\n")
	os.Exit(1)
}

// Exitln is a shim
//go:noinline
func (k *Klogger) Exitln(args ...interface{}) {
	s := fmt.Sprint(args...)
	k.sugar.Error(s, "\n")
	os.Exit(1)
}

// Exitf is a shim
//go:noinline
func Exitf(format string, args ...interface{}) {
	klogger.sugar.Errorf(format, args...)
	os.Exit(1)
}

// Exitf is a shim
//go:noinline
func (k *Klogger) Exitf(format string, args ...interface{}) {
	k.sugar.Errorf(format, args...)
	os.Exit(1)
}

// WithAll fills each arg directly without parsing fields and values
// Only valid for exported fields
func WithAll(args ...interface{}) *Klogger {
	return klogger.WithAll(args...)
}

// WithAll fills each arg directly without parsing fields and values
// Only valid for exported fields
func (k *Klogger) WithAll(args ...interface{}) *Klogger {
	newSugar := k.sugar
	for _, arg := range args {
		t := reflect.TypeOf(arg)
		newSugar = newSugar.With(t.Name(), arg)
	}
	return &Klogger{
		sugar: newSugar,
	}
}

// With fills k-v of a struct into a logger, however it's relatively slow
func With(args ...interface{}) *Klogger {
	return klogger.With(args...)
}

// With fills k-v of a struct into a logger, however it's relatively slow
// Only struct and map will be accepted:
//   * struct: only exported field will be added
//   * map: only accept string type as key
func (k *Klogger) With(args ...interface{}) *Klogger {
	newSugar := k.sugar
	for i := 0; i < len(args); i++ {
		arg := args[i]
		t := reflect.TypeOf(arg)
		v := reflect.ValueOf(arg)
		switch t.Kind() {
		case reflect.Struct:
			for i := 0; i < t.NumField(); i++ {
				k := t.Field(i).Name
				if f := v.Field(i); f.CanInterface() {
					newSugar = newSugar.Desugar().With(zap.Any(k, f.Interface())).Sugar()
				}
			}
		case reflect.Map:
			iter := v.MapRange()
			for iter.Next() {
				if key, val := iter.Key(), iter.Value(); key.Kind() == reflect.String && val.CanInterface() {
					newSugar = newSugar.Desugar().With(zap.Any(key.String(), val.Interface())).Sugar()
				}
			}
		default:
			// other types are not supported yet
		}
	}
	return &Klogger{
		sugar: newSugar,
	}
}

// WithFields requires user to fill in k-v pairs
func WithFields(args ...interface{}) *Klogger {
	return klogger.WithFields(args...)
}

// WithFields requires user to fill in k-v pairs
func (k *Klogger) WithFields(args ...interface{}) *Klogger {
	newSugar := k.sugar.With(args...)
	return &Klogger{
		sugar: newSugar,
	}
}
