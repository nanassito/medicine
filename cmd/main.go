package main

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/nanassito/medicine/pkg/handlers"
	"github.com/nanassito/medicine/pkg/models"
	sheets "google.golang.org/api/sheets/v4"
)

const Year = 365 * 24 * time.Hour

func main() {
	sheets.NewService(context.Background())
	r := mux.NewRouter()
	handler := &handlers.MedicineHandler{
		People: map[models.Person]*models.PersonCfg{
			models.Aline: {
				Birth:    time.Date(1988, 01, 28, 0, 0, 0, 0, time.UTC),
				NextDose: map[models.Medicine]time.Time{},
			},
			models.Dorian: {
				Birth:    time.Date(1989, 5, 9, 0, 0, 0, 0, time.UTC),
				NextDose: map[models.Medicine]time.Time{},
			},
			models.Zaya: {
				Birth:    time.Date(2021, 2, 7, 0, 0, 0, 0, time.UTC),
				NextDose: map[models.Medicine]time.Time{},
			},
			models.Azel: {
				Birth:    time.Date(2023, 6, 27, 0, 0, 0, 0, time.UTC),
				NextDose: map[models.Medicine]time.Time{},
			},
		},
		Medicine: map[models.Medicine]*models.MedicineCfg{
			models.ChildrenIbuprofen: {
				Posology: []models.PosologyEntry{
					{
						OlderThan: 0 * Year,
						Interval:  8 * time.Hour,
						Quantity:  "3ml",
					},
					{
						OlderThan: 2 * Year,
						Interval:  8 * time.Hour,
						Quantity:  "5ml",
					},
					{
						OlderThan: 4 * Year,
						Interval:  8 * time.Hour,
						Quantity:  "7.5ml",
					},
					{
						OlderThan: 6 * Year,
						Interval:  8 * time.Hour,
						Quantity:  "10ml",
					},
					{
						OlderThan: 9 * Year,
						Interval:  8 * time.Hour,
						Quantity:  "12.5ml",
					},
					{
						OlderThan: 11 * Year,
						Interval:  8 * time.Hour,
						Quantity:  "15ml",
					},
				},
			},
			models.InfantAcetaminophen: {
				Posology: []models.PosologyEntry{
					{
						OlderThan: 0 * Year,
						Interval:  6 * time.Hour,
						Quantity:  "3ml",
					},
					{
						OlderThan: 2 * Year,
						Interval:  6 * time.Hour,
						Quantity:  "5ml",
					},
				},
			},
		},
	}
	slog.Info("config", "handler", handler)
	handler.Register(r)

	http.Handle("/", r)
	http.ListenAndServe(":8080", nil)
}
