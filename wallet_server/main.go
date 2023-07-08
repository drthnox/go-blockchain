package main

import (
	"flag"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
	"log"
	"path/filepath"
	"strconv"
	"strings"
)

var config viper.Viper
var env string

func init() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
}

func main() {
	initConfig()
	initLogging()
	port, gateway := initGateway()
	app := NewWalletServer(uint16(port), *gateway)
	app.Run()
}

func initGateway() (int, *string) {
	port := *flag.Int("port", config.GetInt("port"), "TCP Port Number for Wallet Server")
	s := "http://127.0.0.1:" + strconv.Itoa(port)
	gateway := flag.String("gateway", "http://127.0.0.1:"+s, "Blockchain Gateway")
	flag.Parse()
	return port, gateway
}

func initLogging() {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	env := config.GetString("env")
	if strings.ToLower(env) == "dev" {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
}

func initConfig() {
	absPath, err := filepath.Abs("./wallet_server")
	if err != nil {
		panic(err)
	}
	config := viper.New()
	file := filepath.Join(absPath, ".env")
	config.SetConfigFile(file)
	config.SetConfigType("json")
	if err := config.ReadInConfig(); err != nil {
		log.Fatalf("Error to reading config file, %s", err)
	}
}
