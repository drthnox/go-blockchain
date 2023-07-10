package wallet

type BlockchainServerConfig struct {
	Host string `json:"host"`
	Port uint16 `json:"port"`
}

type ServerConfig struct {
	Host              string                    `json:"host"`
	Port              uint16                    `json:"port"`
	BlockchainServers []*BlockchainServerConfig `json:"blockchain_servers"`
}
