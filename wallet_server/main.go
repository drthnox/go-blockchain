package main

import (
	"flag"
	"fmt"
	log "github.com/sirupsen/logrus"
	"go-blockchain/utils"
	"net/url"

	//"github.com/rs/zerolog"
	"github.com/spf13/viper"
	"path/filepath"
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
	host := GetHost()
	port := GetPort()
	gw := fmt.Sprintf("http://%s:%d", *host, port)
	gateway := *flag.String("gateway", gw, "Blockchain Gateway")
	// if gateway is valid, split host:port
	_, err := url.ParseRequestURI(gateway)
	if err != nil {
		panic(err)
	}
	return uint16(port), *host
}

func GetPort() int {
	port := *flag.Int("port", config.GetInt("port"), "TCP Port Number for Wallet Server")
	flag.Parse()
	if port == 0 {
		port = config.GetInt("port")
		if port == 0 {
			log.Infof("No port defined in env - using 0")
			port = 9000
		}
	}
	return port
}

func GetHost() *string {
	host := *flag.String("host", config.GetString("host"), "TCP IP Address for Wallet Server")
	flag.Parse()
	if !utils.CheckIPAddress(host) {
		log.Errorf("ERROR: host IP invalid: %s - looking for env setting", host)
		host = config.GetString("host")
		if !utils.CheckIPAddress(host) {
			log.Errorf("ERROR: config host IP invalid: %s", host)
			log.Info("Falling back to default host: 0.0.0.0")
			host = "0.0.0.0"
		}
	}
	log.Infof("Using host %s", host)
	return &host
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
