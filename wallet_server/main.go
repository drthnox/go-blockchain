package main

import (
	"flag"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"go-blockchain/wallet"
	"net/url"
	"path/filepath"
	"strings"
)

var config viper.Viper
var env string

type Gateway struct {
	host string
	port uint16
}

func main() {
	initConfig()
	initLogging()
	host, port := initWalletServer()
	gateway := GetBlockchainGateway()
	NewWalletServer(host, port, gateway).Run()
}

func initWalletServer() (string, uint16) {
	host := config.GetString("host")
	port := config.GetUint16("port")
	return host, port
}

func GetBlockchainGateway() *Gateway {
	// read in blockchain_servers[] from config
	gatewayHost := config.GetString("blockchain_servers.0.host")
	gatewayPort := config.GetUint16("blockchain_servers.0.port")
	gw := fmt.Sprintf("http://%s:%d", gatewayHost, gatewayPort)
	gateway := *flag.String("gateway", gw, "Blockchain Gateway")
	flag.Parse()
	_, err := url.ParseRequestURI(gateway)
	if err != nil {
		log.Errorf("ERROR: Invalid gateway %v", gateway)
		panic(err)
	}
	return &Gateway{
		host: gatewayHost,
		port: gatewayPort,
	}
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
		log.Fatalf("Error reading config file, %s", err)
	}
	var walletConfig wallet.ServerConfig
	config.Unmarshal(&walletConfig)
}
