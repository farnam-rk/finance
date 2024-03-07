package main

import (
	"company/finance/config"
	"company/finance/db"
	"company/finance/server"
	"company/finance/web"
)

// ConfigFile cofig file path
const ConfigFile = "config/config.json"

// init initiate required items
func init() {

	config.LoadConfig(ConfigFile)
	db.DBconnection()

}

// main main function of the project
func main() {

	web.Register()
	server.Run()

}
