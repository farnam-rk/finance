package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	validator "github.com/asaskevich/govalidator"
)

// ConfigFile cofig file path
const ConfigFile = "config/config.json"

func init() {
	LoadConfig(ConfigFile)
}

// Config for Server
type Config struct {
	DBUsername            string
	DBPassword            string
	Replica               string
	DBName                string `valid:"required"`
	DBHost1               string `valid:"required"`
	DBPort1               int
	DBHost2               string `valid:"required"`
	DBPort2               int
	DBDriver              string   `valid:"required"`
	ListenerPort          string   `json:"ListenerPort"`
	ConfigTimeExportDaily ConfTime `json:"ConfigTimeExportDaily"`
	Cumulative            int      `json:"Cumulative"`
	JWTSecret             string   `json:"JWTSecret"`
}

// Config for time server
type ConfTime struct {
	Hour int `json:"hour"`
	Min  int `json:"min"`
	Sec  int `json:"sec"`
	NSec int `json:"nsec"`
}

var DefaultConfigTime = ConfTime{
	Hour: 16, Min: 01, Sec: 30, NSec: 00,
}

// DefaultConfig Values for Server
var DefaultConfig = Config{
	DBUsername:            "",
	DBPassword:            "",
	Replica:               "",
	DBName:                "finance",
	DBHost1:               "127.0.0.1",
	DBHost2:               "127.0.0.1",
	DBDriver:              "mongodb",
	DBPort1:               27017,
	DBPort2:               27017,
	ConfigTimeExportDaily: DefaultConfigTime,
	Cumulative:            0,
	JWTSecret:             "wellcom to finance",
}

// ConfigValue Variable
var ConfigValue *Config

// LoadConfig from File or if File doesn't Exist from Deafault Values
func LoadConfig(path string) {
	log.Println("loading config file")
	config := &DefaultConfig
	if _, err := os.Stat(path); err == nil {
		f, err := ioutil.ReadFile(path)
		if err != nil {
			log.Println("Error in Reading Config File" + err.Error())
		}
		if err := json.Unmarshal(f, &config); err != nil {
			log.Println("config error:" + err.Error())
		}
	} else {
		log.Println(err)
	}
	if result, _ := validator.ValidateStruct(config); !result {
		log.Println("Loading Config required faild")
	}
	ConfigValue = config
}
