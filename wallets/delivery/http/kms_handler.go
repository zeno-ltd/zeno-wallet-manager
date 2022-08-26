package http

import (
	"encoding/json"
	"net/http"
	"os"
	"os/exec"

	"github.com/go-chi/chi/v5"
	"github.com/zeno/zeno-wallet-manager/domain"
	render "github.com/zeno/zeno-wallet-manager/http/middleware"
	log "go.uber.org/zap"
)

type kmsHandler struct {
	logger     *log.Logger
	secretOpts *domain.SecretOpts
	kmsOpts    *domain.KmsOpts
}

//SetupKmsHandler endpoints to query kms wallets
func SetupKmsHandler(router *chi.Mux, secretOpts *domain.SecretOpts, kmsOpts *domain.KmsOpts, logger *log.Logger) http.Handler {
	handler := &kmsHandler{
		secretOpts: secretOpts,
		kmsOpts:    kmsOpts,
		logger:     logger,
	}
	kmsRouter := chi.NewRouter()

	kmsRouter.Get("/address", handler.CreateAddress)
	kmsRouter.Post("/signer", handler.FetchSigner)
	return kmsRouter
}

// CreateAddress will create new address that is linked to an account
// tatum kms cli command
func (h *kmsHandler) CreateAddress(w http.ResponseWriter, r *http.Request) {
	walletCfg := &domain.WalletConfig{}
	if err := json.NewDecoder(r.Body).Decode(walletCfg); err != nil {
		render.JSON(w, http.StatusBadRequest, render.Map{"status": "error", "data": render.NewError(http.StatusBadRequest, err.Error())})
		return
	}
	var response []byte
	var cmdErr error
	if walletCfg.Network == "testnet" {
		cmd := exec.Command(h.kmsOpts.NodeExec, h.kmsOpts.KmsCMD, "--testnet", "getaddress", walletCfg.WalletID.String(), walletCfg.Index)
		cmd.Env = os.Environ()
		cmd.Env = append(cmd.Env, "TATUM_KMS_PASSWORD="+h.kmsOpts.KmsPassword)
		response, cmdErr = cmd.CombinedOutput()
	} else {
		cmd := exec.Command(h.kmsOpts.NodeExec, h.kmsOpts.KmsCMD, "getaddress", walletCfg.WalletID.String(), walletCfg.Index)
		cmd.Env = os.Environ()
		cmd.Env = append(cmd.Env, "TATUM_KMS_PASSWORD="+h.kmsOpts.KmsPassword)
		response, cmdErr = cmd.CombinedOutput()
	}
	if cmdErr != nil {
		render.JSON(w, http.StatusInternalServerError, render.Map{"status": "error", "data": render.NewError(http.StatusBadRequest, cmdErr.Error()+": "+string(response))})
		return
	}
	address := &domain.NewAddress{}
	err := json.Unmarshal(response, address)
	if err != nil {
		render.JSON(w, http.StatusInternalServerError, render.Map{"status": "error", "data": render.NewError(http.StatusBadRequest, "KMS: Failed to parse response for getaddress command!")})
		return
	}
	render.JSON(w, http.StatusOK, render.Map{"status": "success", "data": address})
}

// FetchSigner will get the private key for the transactional custodial wallet and encrypt it and
// send it across to the backend, the backend will descrypt it to use to sign transactions.
func (h *kmsHandler) FetchSigner(w http.ResponseWriter, r *http.Request) {
	walletCfg := &domain.WalletConfig{}
	if err := json.NewDecoder(r.Body).Decode(walletCfg); err != nil {
		render.JSON(w, http.StatusBadRequest, render.Map{"status": "error", "data": render.NewError(http.StatusBadRequest, err.Error())})
		return
	}
	var response []byte
	var cmdErr error
	if walletCfg.Network == "testnet" {
		cmd := exec.Command(h.kmsOpts.NodeExec, h.kmsOpts.KmsCMD, "--testnet", "getprivatekey", walletCfg.WalletID.String(), walletCfg.Index)
		cmd.Env = os.Environ()
		cmd.Env = append(cmd.Env, "TATUM_KMS_PASSWORD="+h.kmsOpts.KmsPassword)
		response, cmdErr = cmd.CombinedOutput()
	} else {
		cmd := exec.Command(h.kmsOpts.NodeExec, h.kmsOpts.KmsCMD, "getprivatekey", walletCfg.WalletID.String(), walletCfg.Index)
		cmd.Env = os.Environ()
		cmd.Env = append(cmd.Env, "TATUM_KMS_PASSWORD="+h.kmsOpts.KmsPassword)
		response, cmdErr = cmd.CombinedOutput()
	}
	if cmdErr != nil {
		render.JSON(w, http.StatusInternalServerError, render.Map{"status": "error", "data": render.NewError(http.StatusBadRequest, cmdErr.Error()+": "+string(response))})
		return
	}
	signer := &domain.CustodialSigner{}
	err := json.Unmarshal(response, signer)
	if err != nil {
		render.JSON(w, http.StatusInternalServerError, render.Map{"status": "error", "data": render.NewError(http.StatusBadRequest, "KMS: Failed to parse response for getprivatekey command!")})
		return
	}
	var data map[string]string
	encrypted, err := render.Encrypt(signer.PrivateKey, h.secretOpts.WalletCipherKey)
	if err != nil {
		render.JSON(w, http.StatusInternalServerError, render.Map{"status": "error", "data": render.NewError(http.StatusBadRequest, "KMS: "+err.Error())})
		return
	}
	data = make(map[string]string, 1)
	data["signer"] = encrypted
	render.JSON(w, http.StatusOK, render.Map{"status": "success", "data": data})
}
