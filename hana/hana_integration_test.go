// +build integration

package hana

import (
	"bytes"
	"encoding/gob"
	"errors"
	"testing"
	"time"

	"github.com/intelsdi-x/pulse/control/plugin"
	"github.com/intelsdi-x/pulse/core/ctypes"

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
		config["database"] = ctypes.ConfigValueStr{Value: "PULSE_TEST"}
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
		err = sp.Publish(plugin.PulseGOBContentType, buf.Bytes(), *cfg)
		Convey("So publishing metrics should not result in an error", func() {
			So(err, ShouldBeNil)
		})
	})
}
