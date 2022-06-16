package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/zeno/zeno-wallet-manager/config"
	"github.com/zeno/zeno-wallet-manager/wallets"
)

func main() {

	app := fiber.New()
	api := app.Group("/v1", requestid.New())
	api.Use(logger.New(logger.Config{
		Format: "${pid} ${locals:requestid} ${status} - ${method} ${path}",
	}))
	api.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))
	api.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed,
	}))

	/* routes for wallet functionality*/
	walletGroup := api.Group("/wallets")
	walletGroup.Post("/address", wallets.CreateAddress)
	walletGroup.Post("/signer", wallets.FetchSigner)

	app.Listen(config.Get("HTTP_KMS_PORT"))
}
