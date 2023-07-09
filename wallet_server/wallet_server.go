package main

import (
	"bytes"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"go-blockchain/block"
	"go-blockchain/utils"
	"go-blockchain/wallet"
	"html/template"
	"io"
	"net/http"
	"path"
	"strconv"
)

const tempDir = "wallet_server/templates"

type WalletServer struct {
	port    uint16
	gateway string
}

func NewWalletServer(port uint16, gateway string) *WalletServer {
	return &WalletServer{port, gateway}
}

func (ws *WalletServer) Port() uint16 {
	return ws.port
}

func (ws *WalletServer) Gateway() string {
	return ws.gateway
}

func (ws *WalletServer) Index(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		t, _ := template.ParseFiles(path.Join(tempDir, "index.html"))
		t.Execute(w, "")
	default:
		log.Error("ERROR: Invalid HTTP Method")
	}
}

func (ws *WalletServer) Wallet(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodPost:
		w.Header().Add("Content-Type", "application/json")
		myWallet := wallet.NewWallet()
		m, _ := myWallet.MarshalJSON()
		io.WriteString(w, string(m[:]))
	default:
		log.Printf("ERROR: Invalid HTTP Method")
		w.WriteHeader(http.StatusBadRequest)
	}
}

func (ws *WalletServer) CreateTransaction(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodPost:
		decoder := json.NewDecoder(req.Body)
		var t wallet.TransactionRequest
		err := decoder.Decode(&t)
		if err != nil {
			log.Fatalf("ERROR: %v", err)
			io.WriteString(w, string(utils.JsonStatus("fail")))
			return
		}
		if !t.Validate() {
			log.Fatal("ERROR: Missing fields")
			io.WriteString(w, string(utils.JsonStatus("fail")))
			return
		}

		publicKey := utils.PublicKeyFromString(*t.SenderPublicKey)
		privateKey := utils.PrivateKeyFromString(*t.SenderPrivateKey, publicKey)
		value, err := strconv.ParseFloat(*t.Value, 32)
		if err != nil {
			log.Error("ERROR: parse error")
			io.WriteString(w, string(utils.JsonStatus("fail")))
			return
		}
		value32 := float32(value)
		w.Header().Add("Content-Type", "application/json")
		transaction := wallet.NewTransaction(privateKey, publicKey, *t.SenderBlockchainAddress, *t.RecipientBlockchainAddress, value32)
		signature := transaction.GenerateSignature()
		signatureStr := signature.String()
		bt := &block.TransactionRequest{
			SenderBlockchainAddress:    t.SenderBlockchainAddress,
			RecipientBlockchainAddress: t.RecipientBlockchainAddress,
			SenderPublicKey:            t.SenderPublicKey,
			Value:                      &value32,
			Signature:                  &signatureStr,
		}
		m, err := json.Marshal(bt)
		if err != nil {
			log.Error("ERROR: Cannot marshall block transaction request")
			io.WriteString(w, string(utils.JsonStatus("fail")))
			return
		}
		buf := bytes.NewBuffer(m)
		resp, err := http.Post(ws.Gateway()+"/transactions", "application/json", buf)
		if err != nil {
			log.Error("ERROR: Call to blockchain server via gateway " + ws.Gateway() + " failed")
			io.WriteString(w, string(utils.JsonStatus("fail")))
			return
		}
		if resp.StatusCode != http.StatusCreated {
			io.WriteString(w, string(utils.JsonStatus("fail")))
			return
		}
		io.WriteString(w, string(utils.JsonStatus("success")))

	default:
		w.WriteHeader(http.StatusBadRequest)
		log.Error("ERROR: Invalid HTTP Method")
	}
}

func (ws *WalletServer) Run() {
	server := "0.0.0.0:" + strconv.Itoa(int(ws.Port()))
	log.Infof("Starting Wallet Server at %s", server)
	ws.RegisterHandlers()
	log.Fatal(http.ListenAndServe(server, nil))
}

func (ws *WalletServer) RegisterHandlers() {
	http.HandleFunc("/", ws.Index)
	http.HandleFunc("/wallet", ws.Wallet)
	http.HandleFunc("/transaction", ws.CreateTransaction)
}
