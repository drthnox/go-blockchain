package main

import (
	"flag"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"log"
	"path/filepath"
	"strconv"
	"strings"
)

var sugar *zap.SugaredLogger

func init() {

	logger, _ := zap.NewProduction()
	defer logger.Sync() // flushes buffer, if any
	sugar = logger.Sugar()
	sugar.Named("WalletServer")
}

func main() {
	absPath, err := filepath.Abs("./wallet_server")
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
	sugar.Named("WalletServer")
	port := *flag.Int("port", config.GetInt("port"), "TCP Port Number for Wallet Server")
	s := "http://127.0.0.1:" + strconv.Itoa(port)
	gateway := flag.String("gateway", "http://127.0.0.1:"+s, "Blockchain Gateway")
	flag.Parse()
	app := NewWalletServer(uint16(port), *gateway)
	app.Run()
}
