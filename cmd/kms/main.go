package main

import (
	"context"
	"net/http"
	"time"

	"github.com/caarlos0/env/v6"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httprate"
	"github.com/zeno/zeno-wallet-manager/domain"
	kmiddleware "github.com/zeno/zeno-wallet-manager/http/middleware"
	walletApi "github.com/zeno/zeno-wallet-manager/wallets/delivery/http"
	log "go.uber.org/zap"
)

type config struct {
	KmsPassword     string `env:"TATUM_KMS_PASSWORD,notEmpty,unset"`
	WorkDIR         string `env:"WORKDIR"`
	HTTPKmsPort     string `env:"HTTP_KMS_PORT,notEmpty"`
	NodeExec        string `env:"NODE_EXEC,notEmpty"`
	KmsCMD          string `env:"KMS_CMD,notEmpty,unset"`
	WalletCipherKey string `env:"WALLET_CIPHER_KEY,notEmpty,unset"`
	Deployment      string `env:"ENV,notEmpty"`
}

//App is the main app struct
type App struct {
	router *chi.Mux
	config *config
	logger *log.Logger
}

// ErrorResponse is a custom error response struct for apis
type ErrorResponse struct {
	FailedField string
	Tag         string
	Value       string
}

func main() {

	kms := &App{}
	config := kms.loadConfig()
	logger := kmiddleware.GetLogger(config.Deployment)
	kms.logger = logger
	kms.router = chi.NewRouter()
	kms.config = config
	kms.MountHandlers()
	http.ListenAndServe(kms.config.HTTPKmsPort, kms.router)
}

//MountHandlers mounts all handlers for endpoints
func (kms *App) MountHandlers() {
	kms.router.Use(kmiddleware.RequestID)
	kms.router.Use(kmiddleware.Logger(kmiddleware.GetLogger(kms.config.Deployment)))
	kms.router.Use(middleware.Recoverer)
	kms.router.Use(middleware.RealIP)
	kms.router.Use(middleware.Timeout(time.Second * 30))
	kms.router.Use(middleware.URLFormat)
	kms.router.Use(middleware.AllowContentEncoding("deflate", "gzip"))
	kms.router.Use(middleware.Compress(1))
	kms.router.Use(middleware.AllowContentType("application/json"))
	kms.router.Use(kmiddleware.Heartbeat("/health"))
	kms.router.Use(httprate.LimitByIP(100, 1*time.Minute))
	kms.router.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"https://*", "http://*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type", "X-API-KEY"},
		MaxAge:         300, // Maximum value not ignored by any of major browsers
	}))

	secretOps := &domain.SecretOpts{
		JwtSecret:       "",
		WalletCipherKey: kms.config.WalletCipherKey,
	}

	kmsOpts := &domain.KmsOpts{
		NodeExec: kms.config.NodeExec,
		KmsCMD:   kms.config.KmsCMD,
	}

	walletRouter := walletApi.SetupKmsHandler(kms.router, secretOps, kmsOpts, kms.logger)
	kms.router.Route("/v1", func(r chi.Router) {
		r.Use(kms.apiVersionCtx("v1"))
		r.Mount("/wallets", walletRouter)
	})

	kms.router.NotFound(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
		kmiddleware.JSON(w, http.StatusBadRequest, kmiddleware.Map{"status": "error", "data": kmiddleware.NewError(http.StatusBadRequest, "route does not exist")})

	})
	kms.router.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(405)
		kmiddleware.JSON(w, http.StatusBadRequest, kmiddleware.Map{"status": "error", "data": kmiddleware.NewError(http.StatusBadRequest, "method is not valid")})
	})
}

//LoadConfig loads application env config in the Config struct
func (kms *App) loadConfig() *config {
	zConfig := &config{}
	if err := env.Parse(zConfig); err != nil {
		panic("Environment variable not set!" + err.Error())
	}
	return zConfig
}

func (kms *App) apiVersionCtx(version string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r = r.WithContext(context.WithValue(r.Context(), kmiddleware.ZenAPIVersion{}, version))
			next.ServeHTTP(w, r)
		})
	}
}
