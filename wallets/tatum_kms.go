package wallets

import (
	"encoding/json"
	"os/exec"
	"strconv"

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
		response, cmdErr = exec.Command(config.Get("NODE_EXEC"), config.Get("TATUM_KMS"), "--testnet", "getaddress", walletCfg.WalletID.String(), strconv.Itoa(walletCfg.Index)).Output()
	} else {
		response, cmdErr = exec.Command(config.Get("NODE_EXEC"), config.Get("TATUM_KMS"), "getaddress", walletCfg.WalletID.String(), strconv.Itoa(walletCfg.Index)).Output()
	}
	if cmdErr != nil {
		return ctx.JSON(fiber.Map{"status": "error", "data": fiber.NewError(fiber.StatusInternalServerError, cmdErr.Error())})
	}
	address := &NewAddress{}
	err := json.Unmarshal(response, address)
	if err != nil {
		return ctx.JSON(fiber.Map{"status": "error", "data": fiber.NewError(fiber.StatusInternalServerError, "TATUM KMS: Failed to parse response for getaddress command!")})
	}
	return ctx.JSON(fiber.Map{"status": "success", "data": address})
}
