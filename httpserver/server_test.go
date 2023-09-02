package httpserver

import (
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/goinbox/golog"
	"github.com/goinbox/pcontext"
)

func TestServer(t *testing.T) {
	w, _ := golog.NewFileWriter("/tmp/test.log", 0)
	logger := golog.NewSimpleLogger(w, golog.NewSimpleFormater()).
		SetLogLevel(golog.LevelDebug).
		With(&golog.Field{
			Key:   "pid",
			Value: os.Getpid(),
		})
	ctx := pcontext.NewSimpleContext(logger)

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		logger.Info("start sleep")
		time.Sleep(time.Second * 10)
		logger.Info("end sleep")
	})

	server := NewServer("127.0.0.1:8081", mux)
	err := server.ListenAndServe(ctx)
	if err != nil {
		logger.Error("server.ListenAndServe error", golog.ErrorField(err))
	}
}
