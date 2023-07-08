package main

import (
	"flag"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"log"
	"path/filepath"
	"strings"
)

var sugar *zap.SugaredLogger
var config *viper.Viper

func init() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
}

func main() {
	absPath, err := filepath.Abs("./blockchain_server")
	if err != nil {
		panic(err)
	}
	file := filepath.Join(absPath, ".env")
	config := viper.New()
	config.SetConfigFile(file)
	config.SetConfigType("json")
	env := config.GetString("env")
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if strings.ToLower(env) == "dev" {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
	if err := config.ReadInConfig(); err != nil {
		log.Fatalf("Error to reading config file, %s", err)
	}
	port := flag.Uint("port", config.GetUint("port"), "TCP Port Number for Blockchain Server")
	flag.Parse()
	app := NewBlockchainServer(uint16(*port))
	app.Run()
}
