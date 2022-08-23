package wallets

import (
	"encoding/json"
	"os"
	"os/exec"

	"github.com/gofiber/fiber/v2"
	"github.com/zeno/zeno-wallet-manager/config"
)

// CreateAddress will create new address that is linked to an account
// tatum kms cli command
func CreateAddress(ctx *fiber.Ctx) error {
	walletCfg := &WalletConfig{}
	if err := ctx.BodyParser(&walletCfg); err != nil {
		return ctx.JSON(fiber.Map{"status": "error", "data": fiber.NewError(fiber.StatusBadRequest, err.Error())})
	}
	var response []byte
	var cmdErr error
	if walletCfg.Network == "testnet" {
		cmd := exec.Command(config.KmsConfig.NodeExec, config.KmsConfig.KmsCMD, "--testnet", "getaddress", walletCfg.WalletID.String(), walletCfg.Index)
		cmd.Env = os.Environ()
		cmd.Env = append(cmd.Env, "TATUM_KMS_PASSWORD="+config.KmsConfig.KmsPassword)
		response, cmdErr = cmd.CombinedOutput()
	} else {
		cmd := exec.Command(config.KmsConfig.NodeExec, config.KmsConfig.KmsCMD, "getaddress", walletCfg.WalletID.String(), walletCfg.Index)
		cmd.Env = os.Environ()
		cmd.Env = append(cmd.Env, "TATUM_KMS_PASSWORD="+config.KmsConfig.KmsPassword)
		response, cmdErr = cmd.CombinedOutput()
	}
	if cmdErr != nil {
		return ctx.JSON(fiber.Map{"status": "error", "data": fiber.NewError(fiber.StatusInternalServerError, cmdErr.Error()+": "+string(response))})
	}
	address := &NewAddress{}
	err := json.Unmarshal(response, address)
	if err != nil {
		return ctx.JSON(fiber.Map{"status": "error", "data": fiber.NewError(fiber.StatusInternalServerError, "KMS: Failed to parse response for getaddress command!")})
	}
	return ctx.JSON(fiber.Map{"status": "success", "data": address})
}

// FetchSigner will get the private key for the transactional custodial wallet and encrypt it and
// send it across to the backend, the backend will descrypt it to use to sign transactions.
func FetchSigner(ctx *fiber.Ctx) error {
	walletCfg := &WalletConfig{}
	if err := ctx.BodyParser(&walletCfg); err != nil {
		return ctx.JSON(fiber.Map{"status": "error", "data": fiber.NewError(fiber.StatusBadRequest, err.Error())})
	}
	var response []byte
	var cmdErr error
	if walletCfg.Network == "testnet" {
		cmd := exec.Command(config.KmsConfig.NodeExec, config.KmsConfig.KmsCMD, "--testnet", "getprivatekey", walletCfg.WalletID.String(), walletCfg.Index)
		cmd.Env = os.Environ()
		cmd.Env = append(cmd.Env, "TATUM_KMS_PASSWORD="+config.KmsConfig.KmsPassword)
		response, cmdErr = cmd.CombinedOutput()
	} else {
		cmd := exec.Command(config.KmsConfig.NodeExec, config.KmsConfig.KmsCMD, "getprivatekey", walletCfg.WalletID.String(), walletCfg.Index)
		cmd.Env = os.Environ()
		cmd.Env = append(cmd.Env, "TATUM_KMS_PASSWORD="+config.KmsConfig.KmsPassword)
		response, cmdErr = cmd.CombinedOutput()
	}
	if cmdErr != nil {
		return ctx.JSON(fiber.Map{"status": "error", "data": fiber.NewError(fiber.StatusInternalServerError, cmdErr.Error()+": "+string(response))})
	}
	signer := &CustodialSigner{}
	err := json.Unmarshal(response, signer)
	if err != nil {
		return ctx.JSON(fiber.Map{"status": "error", "data": fiber.NewError(fiber.StatusInternalServerError, "KMS: Failed to parse response for getprivatekey command!")})
	}
	var data map[string]string
	encrypted, err := Encrypt(signer.PrivateKey)
	if err != nil {
		return ctx.JSON(fiber.Map{"status": "error", "data": fiber.NewError(fiber.StatusInternalServerError, "KMS: "+err.Error())})
	}
	data = make(map[string]string, 1)
	data["signer"] = encrypted
	return ctx.JSON(fiber.Map{"status": "success", "data": data})
}
