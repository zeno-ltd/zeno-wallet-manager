package config

import (
	"os"

	"github.com/caarlos0/env/v6"
	log "github.com/sirupsen/logrus"
)

//Config will store all env variables in the struct
type Config struct {
	KmsPassword     string `env:"KMS_PASSWORD,notEmpty,unset"`
	WorkDIR         string `env:"WORKDIR"`
	HTTPKmsPort     string `env:"HTTP_KMS_PORT,notEmpty"`
	NodeExec        string `env:"NODE_EXEC,notEmpty"`
	KmsCMD          string `env:"KMS_CMD,notEmpty,unset"`
	WalletCipherKey string `env:"WALLET_CIPHER_KEY,notEmpty,unset"`
}

//KmsConfig global config object
var KmsConfig *Config

//LoadConfig creates the KmsConfig object on startup
func LoadConfig() {
	KmsConfig = &Config{}
	if err := env.Parse(KmsConfig); err != nil {
		log.Error("Environment variable not set!", err.Error())
		os.Exit(1)
	}
}
