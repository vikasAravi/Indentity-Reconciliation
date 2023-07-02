package config

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"strconv"
	"strings"
	"time"
)

type Configuration struct {
	port         int64
	logLevel     string
	logLocation  string
	dbConfig     DbConfig
	queryTimeout int64
}

type DbConfig struct {
	Host            string
	Port            int64
	Name            string
	User            string
	Password        string
	MaxIdleConn     int64
	ConnMaxIdleTime int64
	MaxOpenConn     int64
}

var config Configuration

type Error struct {
	Error error
}

func panicIfError(err error) {
	if err != nil {
		panic(fmt.Errorf("unable to load config: %v", err))
	}
}

func checkKey(key string) {
	if !viper.IsSet(key) {
		panicIfError(fmt.Errorf("%s key is not set", key))
	}
}

func panicIfErrorForKey(err error, key string) {
	if err != nil {
		panicIfError(fmt.Errorf("could not parse key: %s. Error: %v", key, err))
	}
}

func getIntOrPanic(key string) int64 {
	checkKey(key)
	v, err := strconv.Atoi(viper.GetString(key))
	panicIfErrorForKey(err, key)
	return int64(v)
}

func getStringOrPanic(key string) string {
	checkKey(key)
	return viper.GetString(key)
}

func getBoolOrPanic(key string) bool {
	if !viper.IsSet(key) {
		return false
	}

	v, err := strconv.ParseBool(viper.GetString(key))
	panicIfErrorForKey(err, key)
	return v
}

func SetConfigFileFromArgs(commandArgs []string) {
	if len(commandArgs) > 2 {
		viper.SetConfigName(commandArgs[2])
	} else {
		logrus.Fatal("Failed to startup server, please provide a config file name")
	}
}

func SetConfigFileFromFilePath(filepath string) {
	viper.SetConfigName(filepath)
}

func getStringSliceOrPanic(key string) []string {
	checkKey(key)
	return viper.GetStringSlice(key)
}

func getStringArray(key string) []string {
	stringArray := strings.Split(viper.GetString(key), ",")
	for i, str := range stringArray {
		stringArray[i] = strings.TrimSpace(str)
	}
	return stringArray
}

func Load() {
	viper.SetDefault("log_level", "debug")
	viper.AutomaticEnv()

	viper.AddConfigPath("./")
	viper.AddConfigPath("../")
	viper.AddConfigPath("./profiles")
	viper.SetConfigType("yaml")

	err := viper.ReadInConfig()
	if err != nil {
		return
	}

	config = Configuration{
		port:         getIntOrPanic("port"),
		logLevel:     getStringOrPanic("log_level"),
		logLocation:  getStringOrPanic("log_location"),
		queryTimeout: getIntOrPanic("query_timeout"),
		dbConfig: DbConfig{
			Host:            getStringOrPanic("db_host"),
			Port:            getIntOrPanic("db_port"),
			Name:            getStringOrPanic("db_name"),
			User:            getStringOrPanic("db_user"),
			Password:        getStringOrPanic("db_password"),
			MaxIdleConn:     getIntOrPanic("db_max_idle_conn"),
			ConnMaxIdleTime: getIntOrPanic("db_conn_max_idle_time"),
			MaxOpenConn:     getIntOrPanic("db_max_open_conn"),
		},
	}
}

func Port() int64 {
	return config.port
}

func LogLevel() string {
	return config.logLevel
}

func LogLocation() string {
	return config.logLocation
}

func QueryTimeout() time.Duration {
	return time.Duration(config.queryTimeout) * time.Millisecond
}

func DbConf() DbConfig {
	return config.dbConfig
}
