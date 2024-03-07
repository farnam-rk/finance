package server

import (
	"company/finance/acsslog"
	"company/finance/acsslog/dailyrotate"
	"company/finance/config"
	"fmt"
	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// EP is New Echo Instance
var EP = echo.New()

// LogFile Log File Path
const LogFile = "log/finance.log"

var middlewares = []echo.MiddlewareFunc{}

// init Set Middleware Logger
func init() {
	file, err := acsslog.CreateLogFile(LogFile)
	if err != nil {
		panic(err)
	}
	log.SetOutput(file)
	loggercfg := middleware.DefaultLoggerConfig
	loggercfg.CustomTimeFormat = "2006/01/02 15:04:05"
	loggercfg.Format = `${time_custom} Status: "${status}" Error: "${error}"  URI: "${uri}"  UserAgent: "${user_agent}" Host: "${host}" RemoteIP: "${remote_ip}"` + "\n"
	loggercfg.Output = file
	log := middleware.LoggerWithConfig(loggercfg)
	middlewares = append(middlewares, log)
}

// Run Starts Echo Instance in Specified Port
func Run() {

	println("Version 2024-03-06")

	port := fmt.Sprintf(":%s", config.ConfigValue.ListenerPort)
	EP.Start(port)

}

// LoadThread Load All Needed Thread
func LoadThread() {

}

// CreateAccessFile create accesslog file
func CreateAccessFile(mypath, format string) *dailyrotate.File {

	f, err := dailyrotate.NewFile(mypath+"/2006-01-02."+format, func(path string, didRotate bool) { path = mypath; didRotate = true })
	if err != nil {
		log.Println("error in creating access log file : ", err)
	}

	return f

}
