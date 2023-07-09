package main

import (
	"flag"
	log "github.com/sirupsen/logrus"
	"go-blockchain/utils"

	//"github.com/rs/zerolog"
	"github.com/spf13/viper"
	"path/filepath"
	"strconv"
	"strings"
)

var config viper.Viper
var env string

func main() {
	initConfig()
	initLogging()
	NewWalletServer(initGateway()).Run()
}

func initGateway() (uint16, string) {
	host := *flag.String("host", config.GetString("host"), "TCP IP Address for Wallet Server")
	if !utils.CheckIPAddress(host) {
		log.Errorf("ERROR: host IP invalid: %s - looking for env setting", host)
		host = config.GetString("host")
		if !utils.CheckIPAddress(host) {
			log.Errorf("ERROR: config host IP invalid: %s", host)
			log.Info("Falling back to default host: 0.0.0.0")
			host = "0.0.0.0"
		}
	}
	port := *flag.Int("port", config.GetInt("port"), "TCP Port Number for Wallet Server")
	s := "http://127.0.0.1:" + strconv.Itoa(port)
	gateway := *flag.String("gateway", "http://127.0.0.1:"+s, "Blockchain Gateway")
	flag.Parse()
	return uint16(port), gateway
}

func initLogging() {
	log.SetLevel(log.InfoLevel)
	env := config.GetString("env")
	if strings.ToLower(env) == "dev" {
		log.SetLevel(log.DebugLevel)
	}
}

func initConfig() {
	absPath, err := filepath.Abs("./wallet_server")
	if err != nil {
		panic(err)
	}
	config = *viper.New()
	file := filepath.Join(absPath, ".env")
	config.SetConfigFile(file)
	config.SetConfigType("json")
	if err := config.ReadInConfig(); err != nil {
		log.Fatalf("Error to reading config file, %s", err)
	}
}
