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
	c.v = 1 // enable DEBUG level
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

func TestRobust(t *testing.T) {
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
