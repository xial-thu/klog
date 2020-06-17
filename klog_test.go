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
	"fmt"
	"testing"
)

func TestProduction(t *testing.T) {
	InitFlags(nil)
	klogger.config.v = 1 // enable DEBUG level
	Singleton()

	arg := fmt.Errorf("hello")
	arg2 := fmt.Errorf("world")
	Error(arg)
	Errorf("%s", arg)
	Errorln(arg, arg2)
	ErrorDepth(1, arg)
	Warningf("%s", arg)
	Warning(arg)
	Warningln(arg, arg2)
	WarningDepth(1, arg)
	Infof("%s", arg)
	Infoln(arg, arg2)
	InfoDepth(1, arg)
	Info(arg)
	V(1).Info(arg)
	V(1).Infoln(arg)
	V(1).Infof("%s", arg)
	V(2).Info(arg)
	V(2).Infoln(arg)
	V(2).Infof("%s", arg)
}

func TestWith(t *testing.T) {
	Singleton()

	type S struct {
		A int
		B string
	}
	type ID int64
	type W struct {
		C S
		D ID
	}
	type Q struct {
		D ID
	}

	c := "hello"
	s := S{
		A: 10,
		B: "abc",
	}
	w := W{
		C: s,
		D: ID(1),
	}
	q := Q{
		D: ID(1),
	}

	l1 := WithFields("B", "abc")
	l1.Info("hello")
	WithFields("A", 10, "B", "abc").Info(c) // "A":10,"B":"abc"
	With(s).Info(c)                         // "A":10,"B":"abc"
	WithAll(s).Info(c)                      // "S":{"A":10,"B":"abc"}

	// why not split into args
	With(w).Info(c)       // "C":{"A":10,"B":"abc"},"D":1}
	WithAll(w).Info(c)    // "W":{"C":{"A":10,"B":"abc"},"D":1}
	With(s, q).Info(c)    // "A":10,"B":"abc","D":1
	WithAll(s, q).Info(c) // "S":{"A":10,"B":"abc"},"Q":{"D":1}

	// anomony
	With(struct {
		A int
		B int
	}{1, 2}).Info(c) // "A":1,"B":2
	WithAll(struct {
		A int
		B int
	}{1, 2}).Info(c) // "":{"A":1,"B":2}

	type Y map[string]string
	y := Y{"a": "b", "c": "d"}
	With(y, c, map[struct{ A string }]int{{A: "a"}: 1, {A: "b"}: 2}).Info(c)
	WithAll(y, c, map[struct{ A string }]int{{A: "a"}: 1, {A: "b"}: 2}).Info(c)
}

func TestNoOps(t *testing.T) {
	arg := fmt.Errorf("hello")
	arg2 := fmt.Errorf("world")
	Error(arg)
	Errorf("%s", arg)
	Errorln(arg, arg2)
	ErrorDepth(1, arg)
	Warningf("%s", arg)
	Warning(arg)
	Warningln(arg, arg2)
	WarningDepth(1, arg)
	Infof("%s", arg)
	Infoln(arg, arg2)
	InfoDepth(1, arg)
	Info(arg)
	V(1).Info(arg)
	V(1).Infoln(arg)
	V(1).Infof("%s", arg)
	V(2).Info(arg)
	V(2).Infoln(arg)
	V(2).Infof("%s", arg)
}

func TestUpdateLevel(t *testing.T) {
	Singleton()
	Infof("should-always-print")
	V(1).Infof("should-not-print")
	SetLevel(2)
	SetLevel(5)
	SetLevel(-1)
	V(1).Infof("should-print")
}

func BenchmarkWith(b *testing.B) {
	Singleton()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		With(struct {
			ID   string
			Name string
		}{"0001", "hello"}).Info("world")
	}
}

func BenchmarkWithFields(b *testing.B) {
	Singleton()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		WithFields("ID", "0001", "Name", "hello").Info("world")
	}
}

func BenchmarkWithAll(b *testing.B) {
	Singleton()
	type s struct {
		ID   string
		Name string
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		WithAll(s{"0001", "hello"}).Info("world")
	}
}

func BenchmarkContextWith(b *testing.B) {
	Singleton()
	type s struct {
		ID   string
		Name string
	}
	b.ResetTimer()
	newLogger := With(struct {
		ID   string
		Name string
	}{"0001", "hello"})
	for i := 0; i < b.N; i++ {
		newLogger.Info("world")
	}
}
