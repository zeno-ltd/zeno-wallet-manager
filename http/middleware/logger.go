package middleware

import (
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/goccy/go-json"

	"go.uber.org/zap"
)

//GetLogger instance for each repo, service and handler
func GetLogger(deployment string) *zap.Logger {
	var cfg zap.Config
	var rawJSON []byte
	if deployment == "prod" {
		rawJSON = getProductionConfig()
	} else {
		rawJSON = getDevelopmentConfig()
	}
	if err := json.Unmarshal(rawJSON, &cfg); err != nil {
		panic(err)
	}
	logger := zap.Must(cfg.Build())
	defer logger.Sync()
	return logger
}

func getProductionConfig() []byte {
	return []byte(`{
		"level": "error",
		"encoding": "json",
		"outputPaths": ["stdout", "zeno.logs"],
		"errorOutputPaths": ["stderr"],
		"encoderConfig": {
		  "messageKey": "message",
		  "levelKey": "level",
		  "levelEncoder": "lowercase"
		}
	  }`)
}

func getDevelopmentConfig() []byte {
	return []byte(`{
		"level": "debug",
		"encoding": "json",
		"outputPaths": ["stdout", "/tmp/zeno.logs"],
		"errorOutputPaths": ["stderr"],
		"encoderConfig": {
		  "messageKey": "message",
		  "levelKey": "level",
		  "levelEncoder": "lowercase"
		}
	  }`)
}

// Logger is a middleware that logs the start and end of each request, along
// with some useful data about what was requested, what the response status was,
// and how long it took to return.
func Logger(l *zap.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			t1 := time.Now()
			defer func() {
				duration := time.Duration(time.Since(t1).Milliseconds())
				l.Info("http",
					zap.String("method", r.Method),
					zap.String("path", r.URL.Path),
					zap.String("latency", strconv.Itoa(int(duration))+"ms"),
					zap.Int("status", ww.Status()),
					zap.Int("size", ww.BytesWritten()),
					zap.String("requestid", GetReqID(r.Context())))
			}()

			next.ServeHTTP(ww, r)
		}
		return http.HandlerFunc(fn)
	}
}
