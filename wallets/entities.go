package wallets

import "github.com/google/uuid"

//WalletConfig struct will store the request sent to generate an address
type WalletConfig struct {
	Index    string    `json:"index"`
	Network  string    `json:"network"`
	WalletID uuid.UUID `json:"wallet_id"`
}

//NewAddress maps to the response from cli command for getaddress
type NewAddress struct {
	Address string `json:"address"`
}

//CustodialSigner maps to the response from cli command for getprivatekey
type CustodialSigner struct {
	PrivateKey string `json:"privateKey"`
}
