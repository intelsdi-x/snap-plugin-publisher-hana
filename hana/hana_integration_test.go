// +build integration

/*
http://www.apache.org/licenses/LICENSE-2.0.txt

Copyright 2015 Intel Corporation

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
	"bytes"
	"encoding/gob"
	"errors"
	"testing"
	"time"

	"github.com/intelsdi-x/snap/control/plugin"
	"github.com/intelsdi-x/snap/core/ctypes"

	. "github.com/smartystreets/goconvey/convey"
)

func TestHANAPublish(t *testing.T) {
	var buf bytes.Buffer
	metrics := []plugin.PluginMetricType{
		*plugin.NewPluginMetricType([]string{"test_string"}, time.Now(), "", "example_string"),
		*plugin.NewPluginMetricType([]string{"test_int"}, time.Now(), "", 1),
		*plugin.NewPluginMetricType([]string{"test_string_slice"}, time.Now(), "", []string{"str1", "str2"}),
		*plugin.NewPluginMetricType([]string{"test_string_slice"}, time.Now(), "", []int{1, 2}),
		*plugin.NewPluginMetricType([]string{"test_uint8"}, time.Now(), "", uint8(1)),
	}
	config := make(map[string]ctypes.ConfigValue)
	enc := gob.NewEncoder(&buf)
	enc.Encode(metrics)

	Convey("TestHANAPublish", t, func() {
		config["username"] = ctypes.ConfigValueStr{Value: "root"}
		config["password"] = ctypes.ConfigValueStr{Value: "root"}
		config["database"] = ctypes.ConfigValueStr{Value: "SNAP_TEST"}
		config["host"] = ctypes.ConfigValueStr{Value: "localhost"}
		config["port"] = ctypes.ConfigValueStr{Value: "1433"}
		config["table name"] = ctypes.ConfigValueStr{Value: "info"}
		sp := NewHANAPublisher()
		cp, _ := sp.GetConfigPolicy()
		cfg, _ := cp.Get([]string{""}).Process(config)
		err := sp.Publish("", buf.Bytes(), *cfg)
		Convey("So not passing in a content type should result in an error", func() {
			So(err, ShouldResemble, errors.New("Unknown content type ''"))
		})
		err = sp.Publish(plugin.SnapGOBContentType, buf.Bytes(), *cfg)
		Convey("So publishing metrics should not result in an error", func() {
			So(err, ShouldBeNil)
		})
	})
}
