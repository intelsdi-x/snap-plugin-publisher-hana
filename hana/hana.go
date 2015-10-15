package hana

import (
	"bytes"
	"database/sql"
	"encoding/gob"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	log "github.com/Sirupsen/logrus"

	"github.com/SAP/go-hdb/driver"
	"github.com/intelsdi-x/pulse/control/plugin"
	"github.com/intelsdi-x/pulse/control/plugin/cpolicy"
	"github.com/intelsdi-x/pulse/core/ctypes"
)

const (
	name       = "hana"
	version    = 2
	pluginType = plugin.PublisherPluginType
)

type HANAPublisher struct {
}

func NewHANAPublisher() *HANAPublisher {
	return &HANAPublisher{}
}

// Publish sends data to a HANA server
func (s *HANAPublisher) Publish(contentType string, content []byte, config map[string]ctypes.ConfigValue) error {
	logger := log.New()
	logger.Println("Publishing started")
	var metrics []plugin.PluginMetricType

	switch contentType {
	case plugin.PulseGOBContentType:
		dec := gob.NewDecoder(bytes.NewBuffer(content))
		if err := dec.Decode(&metrics); err != nil {
			logger.Printf("Error decoding: error=%v content=%v", err, content)
			return err
		}
	default:
		logger.Printf("Error unknown content type '%v'", contentType)
		return errors.New(fmt.Sprintf("Unknown content type '%s'", contentType))
	}

	logger.Printf("publishing %v to %v", metrics, config)

	// Open connection and ping to make sure it works
	username := config["username"].(ctypes.ConfigValueStr).Value
	password := config["password"].(ctypes.ConfigValueStr).Value
	host := config["host"].(ctypes.ConfigValueStr).Value
	database := config["database"].(ctypes.ConfigValueStr).Value
	port := config["port"].(ctypes.ConfigValueStr).Value
	tableName := config["tablename"].(ctypes.ConfigValueStr).Value
	tableColumns := "(timestamp VARCHAR(200), source VARCHAR(200), key VARCHAR(200), value VARCHAR(200))"
	db, err := sql.Open(driver.DriverName, "hdb://"+username+":"+password+"@"+host+":"+port+"/"+database)
	defer db.Close()
	if err != nil {
		logger.Printf("Error: %v", err)
		return err
	}
	err = db.Ping()
	if err != nil {
		logger.Printf("Error: %v", err)
		return err
	}

	// Create the table if it's not already there
	_, err = db.Exec("DROP PROCEDURE ifexists")

	if err != nil {
		logger.Printf("Error while dropping procedure: %v", err)
	}

	createTableStr :=
		"CREATE PROCEDURE ifexists( ) LANGUAGE SQLSCRIPT AS myrowid integer;\n" +
			"BEGIN\n" +
			" myrowid := 0;\n" +
			" SELECT COUNT(*) INTO myrowid FROM \"PUBLIC\".\"M_TABLES\" " +
			" WHERE schema_name = '" + database + "' AND table_name = '" + tableName + "';\n" +
			" IF :myrowid = 0 THEN\n" +
			"  exec 'CREATE COLUMN TABLE \"" + database + "\".\"" + tableName + "\" " + tableColumns + "';\n " +
			" END IF;\n" +
			"END;"

	_, err = db.Exec(createTableStr)

	if err != nil {
		logger.Printf("Error while creating procedure: %v", err)
		logger.Printf("Query: %v", createTableStr)
	}

	_, err = db.Exec("CALL ifexists")

	if err != nil {
		logger.Printf("Error while invoking procedure: %v", err)

		return err
	}

	// Put the values into the database with the current time
	tableValues := "VALUES( ?, ?, ?, ? )"
	insert, err := db.Prepare("INSERT INTO " + database + "." + tableName + " " + tableValues)
	if err != nil {
		logger.Printf("Error: %v", err)
		logger.Printf("tablename: " + database + "." + tableName + ", tableValues: " + tableValues)
		return err
	}

	var key, value string
	for _, m := range metrics {
		key = sliceToString(m.Namespace())
		value, err = interfaceToString(m.Data())
		if err == nil {
			_, err := insert.Exec(m.Timestamp(), m.Source(), key, value)
			if err != nil {
				panic(err)
				logger.Printf("Error: %v", err)
			}
		} else {
			logger.Printf("Error: %v", err)
		}
	}

	return err
}

func Meta() *plugin.PluginMeta {
	return plugin.NewPluginMeta(name, version, pluginType, []string{plugin.PulseGOBContentType}, []string{plugin.PulseGOBContentType})
}

func (f *HANAPublisher) GetConfigPolicy() (*cpolicy.ConfigPolicy, error) {
	cp := cpolicy.New()
	config := cpolicy.NewPolicyNode()

	host, err := cpolicy.NewStringRule("host", true, "root")
	handleErr(err)
	host.Description = "HANA host to which we will connect"

	port, err := cpolicy.NewStringRule("port", true, "30017")
	handleErr(err)
	port.Description = "HANA port to which we will connect"

	username, err := cpolicy.NewStringRule("username", true, "root")
	handleErr(err)
	username.Description = "Username to login to the HANA server"

	password, err := cpolicy.NewStringRule("password", true, "root")
	handleErr(err)
	password.Description = "Password to login to the HANA server"

	database, err := cpolicy.NewStringRule("database", true, "PULSE_TEST")
	handleErr(err)
	database.Description = "The HANA database that data will be pushed to"

	tableName, err := cpolicy.NewStringRule("tablename", true, "info")
	handleErr(err)
	tableName.Description = "The HANA table within the database where information will be stored"

	config.Add(
		host,
		username,
		password,
		database,
		port,
		tableName,
	)

	cp.Add([]string{""}, config)
	return cp, nil
}

func handleErr(e error) {
	if e != nil {
		panic(e)
	}
}

func sliceToString(slice []string) string {
	return strings.Join(slice, ", ")
}

// Supported types: []string, []int, int, uint64, float32, float64, string
func interfaceToString(face interface{}) (string, error) {
	var (
		ret string
		err error
	)
	switch val := face.(type) {
	case []string:
		ret = sliceToString(val)
	case []int:
		length := len(val)
		if length == 0 {
			return ret, err
		}
		ret = strconv.Itoa(val[0])
		if length == 1 {
			return ret, err
		}
		for i := 1; i < length; i++ {
			ret += ", "
			ret += strconv.Itoa(val[i])
		}
	case int:
		ret = strconv.Itoa(val)
	case uint64:
		ret = strconv.FormatUint(val, 10)
	case float32:
		ret = strconv.FormatFloat(float64(val), 'E', -1, 32)
	case float64:
		ret = strconv.FormatFloat(val, 'E', -1, 64)
	case string:
		ret = val
	default:
		err = errors.New("Unsupported type: " + reflect.TypeOf(face).String())
	}
	return ret, err
}
