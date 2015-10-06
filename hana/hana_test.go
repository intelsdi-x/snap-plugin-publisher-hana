// +build unit

/*
http://www.apache.org/licenses/LICENSE-2.0.txt


Copyright 2015 Intel Coporation

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

package hana

import (
	"testing"

	"github.com/intelsdi-x/pulse/control/plugin"
	"github.com/intelsdi-x/pulse/control/plugin/cpolicy"
	"github.com/intelsdi-x/pulse/core/ctypes"
	. "github.com/smartystreets/goconvey/convey"
)

func TestHANAPublisherPlugin(t *testing.T) {
	Convey("Meta should return metadata for the plugin", t, func() {
		meta := Meta()
		So(meta.Name, ShouldResemble, name)
		So(meta.Version, ShouldResemble, version)
		So(meta.Type, ShouldResemble, plugin.PublisherPluginType)
	})

	Convey("Create HANA Publisher", t, func() {
		h := NewHANAPublisher()
		Convey("So h should not be nil", func() {
			So(h, ShouldNotBeNil)
		})
		Convey("So h should be of hanaPublisher type", func() {
			So(h, ShouldHaveSameTypeAs, &HANAPublisher{})
		})
		configPolicy, err := h.GetConfigPolicy()
		Convey("h.GetConfigPolicy() should return a config policy", func() {
			Convey("So config policy should not be nil", func() {
				So(configPolicy, ShouldNotBeNil)
			})
			Convey("So we should not get an err retreiving the config policy", func() {
				So(err, ShouldBeNil)
			})
			Convey("So config policy should be a cpolicy.ConfigPolicy", func() {
				So(configPolicy, ShouldHaveSameTypeAs, &cpolicy.ConfigPolicy{})
			})
			testConfig := make(map[string]ctypes.ConfigValue)
			testConfig["host"] = ctypes.ConfigValueStr{Value: "localhost"}
			testConfig["port"] = ctypes.ConfigValueStr{Value: "8086"}
			testConfig["username"] = ctypes.ConfigValueStr{Value: "root"}
			testConfig["password"] = ctypes.ConfigValueStr{Value: "root"}
			testConfig["database"] = ctypes.ConfigValueStr{Value: "test"}
			testConfig["tablename"] = ctypes.ConfigValueStr{Value: "PULSE_TEST"}
			cfg, errs := configPolicy.Get([]string{""}).Process(testConfig)
			Convey("So config policy should process testConfig and return a config", func() {
				So(cfg, ShouldNotBeNil)
			})
			Convey("So testConfig processing should return no errors", func() {
				So(errs.HasErrors(), ShouldBeFalse)
			})
		})
	})
}
