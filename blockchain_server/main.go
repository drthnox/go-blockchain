package main

import (
	"flag"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"log"
	"path/filepath"
	"strings"
)

var sugar *zap.SugaredLogger
var config *viper.Viper

func init() {

}

func main() {
	absPath, err := filepath.Abs("./blockchain_server")
	if err != nil {
		panic(err)
	}

	// Print the absolute path of the file
	log.Printf("=====> %s", absPath)
	config := viper.New()
	file := filepath.Join(absPath, ".env")
	config.SetConfigFile(file)
	config.SetConfigType("json")
	if err := config.ReadInConfig(); err != nil {
		log.Fatalf("Error to reading config file, %s", err)
	}
	env := config.GetString("env")
	logger, _ := zap.NewProduction()
	if strings.ToLower(env) == "dev" {
		logger, _ = zap.NewDevelopment()
	}
	defer logger.Sync() // flushes buffer, if any
	sugar = logger.Sugar()
	sugar.Named("BlockchainServer")
	port := flag.Uint("port", config.GetUint("port"), "TCP Port Number for Blockchain Server")
	flag.Parse()
	app := NewBlockchainServer(uint16(*port))
	app.Run()
}
