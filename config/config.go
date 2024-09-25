package config

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
	"strings"
)

const (
	ServiceName            = "eagle-rest-service"
	serverPort             = "SERVER_PORT"
	applicationEnvLocation = "APP_ENV_LOCATION"
	postgresHost           = "POSTGRES_HOST"
	postgresPort           = "POSTGRES_PORT"
	postgresDb             = "POSTGRES_DB"
	postgresUser           = "POSTGRES_USER"
	postgresPass           = "POSTGRES_PASS"
	logLevel               = "LOG_LEVEL"
	defaultLogLevel        = "info"
	readTimeout            = "SERVER_READ_TIMEOUT"
	writeTimeout           = "SERVER_WRITE_TIMEOUT"
)

const (
	defaultReadTimeout  = 15
	defaultWriteTimeout = 15
)

func LoadConfig(path string) {
	viper.AddConfigPath(path)
	viper.SetConfigName("application")
	viper.SetConfigType("env")
	viper.AutomaticEnv()
	viper.SetDefault(logLevel, defaultLogLevel)
	viper.SetDefault(postgresPort, 5432)
	viper.SetDefault(readTimeout, defaultReadTimeout)
	viper.SetDefault(writeTimeout, defaultWriteTimeout)
	viper.SetDefault(serverPort, "8000")

	err := viper.ReadInConfig()

	if err != nil {
		logrus.Warnf("error while reading config file : %s", err)
	}

}

const missingVariableString = "%s is missing"

func CheckRequiredVariables() error {

	var missing = make([]string, 0)

	if len(GetServerPort()) == 0 {
		missing = append(missing, missingVariableString, serverPort)
	}

	if len(GetPostgresHost()) == 0 {
		missing = append(missing, missingVariableString, postgresHost)
	}

	if len(GetPostgresUser()) == 0 {
		missing = append(missing, missingVariableString, postgresUser)
	}

	if len(GetPostgresPass()) == 0 {
		missing = append(missing, missingVariableString, postgresPass)
	}

	if len(GetPostgresDb()) == 0 {
		missing = append(missing, missingVariableString, postgresDb)
	}

	if len(missing) > 0 {
		return fmt.Errorf("failed to startup the service check following env variables : %s", strings.Join(missing, ", "))
	}
	return nil
}

// getters for all the environment variables

func GetServerPort() string {
	return viper.GetString(serverPort)
}

// for loglevel

func GetLogLevel() logrus.Level {
	level := viper.GetString(logLevel)
	parseLevel, err := logrus.ParseLevel(level)

	if err != nil {
		return 0
	}
	return parseLevel
}

func GetPostgresHost() string {
	return viper.GetString(postgresHost)
}

func GetPostgresUser() string {
	return viper.GetString(postgresUser)
}

func GetPostgresPass() string {
	return viper.GetString(postgresPass)
}

func GetAppEnvLocation() string {
	loc := os.Getenv(applicationEnvLocation)

	if len(loc) == 0 {
		loc = "."
	}
	return loc
}

func GetPostgresDb() string {
	return viper.GetString(postgresDb)
}

func GetPostgresPort() int {
	return viper.GetInt(postgresPort)
}

func GetServiceName() string {
	return viper.GetString(ServiceName)
}

func GetReadTimeout() int {
	return viper.GetInt(readTimeout)
}

func GetWriteTimeout() int {
	return viper.GetInt(writeTimeout)
}
