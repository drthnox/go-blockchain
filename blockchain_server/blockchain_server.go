package main

import (
	log "github.com/sirupsen/logrus"
	"go-blockchain/block"
	"go-blockchain/wallet"
	"io"
	"net/http"
	"strconv"
)

var cache map[string]*block.Blockchain = make(map[string]*block.Blockchain)

type BlockchainServer struct {
	port uint16
}

func NewBlockchainServer(port uint16) *BlockchainServer {
	log.Debugf("Creating Blockchain Server to use port %d", port)
	return &BlockchainServer{port}
}

func (bcs *BlockchainServer) Port() uint16 {
	return bcs.port
}

func (bcs *BlockchainServer) GetBlockchain() *block.Blockchain {
	bc, ok := cache["blockchain"]
	if !ok {
		minersWallet := wallet.NewWallet()
		bc = block.NewBlockchain(minersWallet.BlockchainAddress(), bcs.Port())
		cache["blockchain"] = bc
		log.Infof("private_key %v", minersWallet.PrivateKeyStr())
		log.Infof("publick_key %v", minersWallet.PublicKeyStr())
		log.Infof("blockchain_address %v", minersWallet.BlockchainAddress())
	}
	return bc
}

func (bcs *BlockchainServer) GetChain(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		w.Header().Add("Content-Type", "application/json")
		bc := bcs.GetBlockchain()
		m, _ := bc.MarshalJSON()
		io.WriteString(w, string(m[:]))
	default:
		log.Printf("ERROR: Invalid HTTP Method")
	}
}

func (bcs *BlockchainServer) Run() {
	server := "0.0.0.0:" + strconv.Itoa(int(bcs.Port()))
	log.Infof("Starting Blockchain Server at %s", server)
	bcs.RegisterHandlers()
	log.Fatal(http.ListenAndServe(server, nil))
}

func (bcs *BlockchainServer) RegisterHandlers() {
	http.HandleFunc("/", bcs.GetChain)
}
