package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/nanassito/medicine/pkg/handlers"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

func mustGoogleService() *sheets.Service {
	credentials, err := os.ReadFile("/Users/dorian.jaminaisgrellier/Downloads/medicine-449623-ba166d49b718.json")
	if err != nil {
		log.Fatal("unable to read key file:", err)
	}

	scopes := []string{
		"https://www.googleapis.com/auth/spreadsheets.readonly",
	}
	config, err := google.JWTConfigFromJSON(credentials, scopes...)
	if err != nil {
		log.Fatal("unable to create JWT configuration:", err)
	}

	ctx := context.Background()
	srv, err := sheets.NewService(ctx, option.WithHTTPClient(config.Client(ctx)))
	if err != nil {
		log.Fatalf("unable to retrieve sheets service: %v", err)
	}
	return srv
}

func main() {
	r := mux.NewRouter()
	handler, err := handlers.NewMedicineHandler(mustGoogleService())
	if err != nil {
		log.Fatal("unable to start the service:", err)
	}
	slog.Info("config", "handler", handler)
	handler.Register(r)

	http.Handle("/", r)
	http.ListenAndServe(":8080", nil)
}
