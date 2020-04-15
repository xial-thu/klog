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
	"sync"

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

var (
	logger *zap.Logger
	sugar  *zap.SugaredLogger
	c      config
	once   sync.Once
)

// init as the global no-ops logger so that unit test will not crash
func init() {
	logger = zap.L()
	sugar = zap.S()
}

// Singleton inits an unique logger
func Singleton() {
	once.Do(func() {
		c.zapConfig = zap.NewProductionConfig()

		// change time from ns to formatted
		c.zapConfig.EncoderConfig.TimeKey = "time"
		c.zapConfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

		// due to gaps between zap and klog
		if c.v > 0 {
			c.zapConfig.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
		}
		if !c.alsologtostderr {
			c.zapConfig.OutputPaths = []string{"stdout"}
		}

		var err error
		// trace the real source caller due to munual inline is not supported
		logger, err = c.zapConfig.Build(zap.AddCallerSkip(1))
		if err != nil {
			panic(err)
		}
		sugar = logger.Sugar()
	})
	sugar.Infof("init zap logger...")
}

// InitFlags is a shim
func InitFlags(flagset *flag.FlagSet) {
	flag.IntVar(&c.v, "v", 0, "verbosity of info log")
	flag.BoolVar(&c.alsologtostderr, "alsologtostderr", true, "also write logs to stderr, default to true")
}

// Flush is a shim
func Flush() {
	logger.Sync()
}

// V is a shim
func V(level Level) Verbose {
	return Verbose(int(level) <= c.v)
}

// Info is a shim
//go:noinline
func (v Verbose) Info(args ...interface{}) {
	if bool(v) {
		sugar.Debug(args...)
	}
}

// Infoln is a shim
//go:noinline
func (v Verbose) Infoln(args ...interface{}) {
	if bool(v) {
		s := fmt.Sprint(args...)
		sugar.Debug(s, "\n")
	}
}

// Infof is a shim
//go:noinline
func (v Verbose) Infof(format string, args ...interface{}) {
	if bool(v) {
		sugar.Debugf(format, args...)
	}
}

// Info is a shim
//go:noinline
func Info(args ...interface{}) {
	sugar.Info(args...)
}

// InfoDepth is a shim
//go:noinline
func InfoDepth(depth int, args ...interface{}) {
	sugar.Info(args...)
}

// Infoln is a shim
//go:noinline
func Infoln(args ...interface{}) {
	s := fmt.Sprint(args...)
	sugar.Info(s, "\n")
}

// Infof is a shim
//go:noinline
func Infof(format string, args ...interface{}) {
	logger.Sugar().Infof(format, args...)
}

// Warning is a shim
//go:noinline
func Warning(args ...interface{}) {
	logger.Sugar().Warn(args...)
}

// WarningDepth is a shim
//go:noinline
func WarningDepth(depth int, args ...interface{}) {
	sugar.Warn(args...)
}

// Warningln is a shim
//go:noinline
func Warningln(args ...interface{}) {
	s := fmt.Sprint(args...)
	sugar.Warn(s, "\n")
}

// Warningf is a shim
//go:noinline
func Warningf(format string, args ...interface{}) {
	sugar.Warnf(format, args...)
}

// Error is a shim
//go:noinline
func Error(args ...interface{}) {
	sugar.Error(args...)
}

// ErrorDepth is a shim
//go:noinline
func ErrorDepth(depth int, args ...interface{}) {
	sugar.Error(args...)
}

// Errorln is a shim
//go:noinline
func Errorln(args ...interface{}) {
	s := fmt.Sprint(args...)
	sugar.Error(s, "\n")
}

// Errorf is a shim
//go:noinline
func Errorf(format string, args ...interface{}) {
	sugar.Errorf(format, args...)
}

// Fatal is a shim
//go:noinline
func Fatal(args ...interface{}) {
	sugar.Error(args...)
	os.Exit(255)
}

// FatalDepth is a shim
//go:noinline
func FatalDepth(depth int, args ...interface{}) {
	sugar.Error(args...)
	os.Exit(255)
}

// Fatalln is a shim
//go:noinline
func Fatalln(args ...interface{}) {
	s := fmt.Sprint(args...)
	sugar.Error(s, "\n")
	os.Exit(255)
}

// Fatalf is a shim
//go:noinline
func Fatalf(format string, args ...interface{}) {
	sugar.Errorf(format, args...)
	os.Exit(255)
}

// Exit is a shim
//go:noinline
func Exit(args ...interface{}) {
	sugar.Error(args...)
	os.Exit(1)
}

// ExitDepth is a shim
//go:noinline
func ExitDepth(depth int, args ...interface{}) {
	sugar.Error(args...)
	os.Exit(1)
}

// Exitln is a shim
//go:noinline
func Exitln(args ...interface{}) {
	s := fmt.Sprint(args...)
	sugar.Error(s, "\n")
	os.Exit(1)
}

// Exitf is a shim
//go:noinline
func Exitf(format string, args ...interface{}) {
	logger.Sugar().Errorf(format, args...)
	os.Exit(1)
}
