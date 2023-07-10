package main

import (
	"flag"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"path/filepath"
	"strings"
)

var config viper.Viper

func main() {
	initConfig()
	initLogger()
	port := flag.Uint("port", config.GetUint("port"), "TCP Port Number for Blockchain Server")
	flag.Parse()
	if *port == 0 {
		*port = config.GetUint("port")
	}
	log.Debugf("Configured port: %d", *port)
	NewBlockchainServer(uint16(*port)).Run()
}

func initLogger() {
	env := config.GetString("env")
	log.SetLevel(log.InfoLevel)
	if strings.ToLower(env) == "dev" {
		log.SetLevel(log.DebugLevel)
	}
}

func initConfig() {
	absPath, err := filepath.Abs("./blockchain_server")
	if err != nil {
		panic(err)
	}
	file := filepath.Join(absPath, ".env")
	config = *viper.New()
	config.SetConfigFile(file)
	config.SetConfigType("json")
	if err := config.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}
	log.Printf("=====> %d", config.GetUint("port"))
}
