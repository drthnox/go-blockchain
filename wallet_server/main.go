package main

import (
	"flag"
	log "github.com/sirupsen/logrus"
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
