package main

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"go-blockchain/block"
	"go-blockchain/utils"
	"go-blockchain/wallet"
	"io"
	"net/http"
)

var cache map[string]*block.Blockchain = make(map[string]*block.Blockchain)

type BlockchainServer struct {
	host string
	port uint16
}

func NewBlockchainServer(host string, port uint16) *BlockchainServer {
	log.Debugf("Creating Blockchain Server running at %s:%d", host, port)
	return &BlockchainServer{host, port}
}

func (bcs *BlockchainServer) Host() string {
	return bcs.host
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
		log.Infof("public_key %v", minersWallet.PublicKeyStr())
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
		w.WriteHeader(http.StatusBadRequest)
	}
}

func (bcs *BlockchainServer) Transactions(w http.ResponseWriter, req *http.Request) {
	log.Debug("Handler: Transactions")
	switch req.Method {
	case http.MethodGet:
		w.Header().Add("Content-Type", "application/json")
		bc := bcs.GetBlockchain()
		transactions := bc.TransactionPool()
		m, _ := json.Marshal(struct {
			Transactions []*block.Transaction `json:"transactions"`
			Length       int                  `json:"length"`
		}{
			Transactions: transactions,
			Length:       len(transactions),
		})
		io.WriteString(w, string(m))

	case http.MethodPost:
		decoder := json.NewDecoder(req.Body)
		var t block.TransactionRequest
		err := decoder.Decode(&t)
		if err != nil {
			log.Errorf("ERROR: %v", err)
			io.WriteString(w, string(utils.JsonStatus("fail")))
			return
		}
		if !t.Validate() {
			log.Error("ERROR: Missing fields")
			io.WriteString(w, string(utils.JsonStatus("fail")))
			return
		}
		publicKey := utils.PublicKeyFromString(*t.SenderPublicKey)
		signature := utils.SignatureFromString(*t.Signature)
		bc := bcs.GetBlockchain()
		isCreated := bc.CreateTransaction(*t.SenderBlockchainAddress, *t.RecipientBlockchainAddress, *t.Value, publicKey, signature)
		w.Header().Add("Content-Type", "application/json")
		var m []byte
		if isCreated {
			w.WriteHeader(http.StatusCreated)
			m = utils.JsonStatus("success")
		} else {
			w.WriteHeader(http.StatusBadRequest)
			m = utils.JsonStatus("fail")
		}
		io.WriteString(w, string(m))

	default:
		log.Error("ERROR: Invalid HTTP Method")
		w.WriteHeader(http.StatusBadRequest)
	}
}

func (bcs *BlockchainServer) Run() {
	server := fmt.Sprintf("%s:%d", bcs.Host(), int(bcs.Port()))
	log.Infof("Starting Blockchain Server at %s", server)
	bcs.RegisterHandlers()
	log.Fatal(http.ListenAndServe(server, nil))
}

func (bcs *BlockchainServer) RegisterHandlers() {
	http.HandleFunc("/", bcs.GetChain)
	http.HandleFunc("/transactions", bcs.Transactions)
}
