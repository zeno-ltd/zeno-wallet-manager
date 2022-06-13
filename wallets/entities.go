package wallets

import "github.com/google/uuid"

//WalletConfig struct will store the request sent to generate an address
type WalletConfig struct {
	WalletID uuid.UUID `json:"wallet_id"`
	Index    string    `json:"index"`
	Network  string    `json:"network"`
}

//NewAddress maps to the response from cli command for getaddress
type NewAddress struct {
	Address string `json:"address"`
}
