package config

import (
	"gopkg.in/ini.v1"
	"log"
	"os"
)

type ConfigList struct {
	DriverName string
	DataSourceName string
	Port string
}
var Config ConfigList
func init() {
	config, err := ini.Load("app/config.ini")
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	Config = ConfigList{
		DriverName: config.Section("db").Key("DriverName").String(),
		DataSourceName: config.Section("db").Key("DataSourceName").String(),
		Port: config.Section("server").Key("Port").String(),
	}
}
