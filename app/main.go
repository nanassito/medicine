package main

import (
	"context"
	"flag"
	"log"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gorilla/mux"
	"github.com/nanassito/medicine/pkg/handlers"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

var (
	creds = flag.String("creds", "../creds.json", "Google credential file.")
)

func mustGetCreds() []byte {
	if creds == nil {
		log.Fatal("missing credentials")
	}
	if strings.HasPrefix(*creds, "../") || strings.HasPrefix(*creds, "./") {
		cwd, err := os.Getwd()
		if err != nil {
			log.Fatal("unable to get current working directory:", err)
		}
		*creds = filepath.Join(cwd, *creds)
	}
	credentials, err := os.ReadFile(*creds)
	if err != nil {
		log.Fatal("unable to read key file:", err)
	}
	return credentials
}

func mustGoogleService() *sheets.Service {
	flag.Parse()

	scopes := []string{
		"https://www.googleapis.com/auth/spreadsheets",
	}
	config, err := google.JWTConfigFromJSON(mustGetCreds(), scopes...)
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
	http.ListenAndServe(":80", nil)
}
