package main

import (
	"net/http"

	logruslogger "stout.dev/jisho/internal/adapters/logger"
	jishoclient "stout.dev/jisho/internal/external/jishoClient"
	"stout.dev/jisho/internal/handler"
	"stout.dev/jisho/internal/usecase"
)

func main() {
	mux := http.NewServeMux()
	log := logruslogger.NewLogrusLogger()

	jishoClient := jishoclient.JishoUsecase(log)
	jishoUsecase := usecase.JishoUsecase(jishoClient, log)
	jishoHandler := handler.NewJishoHandler(jishoUsecase, log)

	mux.HandleFunc("/v1/search", jishoHandler.SearchJisho)
	http.ListenAndServe(":8080", mux)
	log.Info("Server ready and listening!", nil)
}
