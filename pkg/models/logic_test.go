package models_test

import (
	"testing"
	"time"

	"github.com/nanassito/medicine/pkg/models"
)

func TestCanTake(t *testing.T) {
	tests := []struct {
		name       string
		snapshot   models.Snapshot
		person     models.Person
		medicine   models.Medicine
		wantResult bool
		wantMsg    string
	}{
		{
			name: "Person too young",
			snapshot: models.Snapshot{
				People: models.PeopleSlice{
					{Name: "John", Birth: time.Now().AddDate(-1, 0, 0)}, // 1 year old
				},
				Medicines: models.MedicinesMap{
					"Aspirin": &models.MedicineCfg{
						Posology: []models.PosologyEntry{
							{OlderThan: 2 * 365 * 24 * time.Hour}, // 2 years
						},
					},
				},
			},
			person:     "John",
			medicine:   "Aspirin",
			wantResult: false,
			wantMsg:    models.ErrTooYoung.Error(),
		},
		{
			name: "Medicine not found",
			snapshot: models.Snapshot{
				People: models.PeopleSlice{
					{Name: "John", Birth: time.Now().AddDate(-10, 0, 0)}, // 10 years old
				},
				Medicines: models.MedicinesMap{},
			},
			person:     "John",
			medicine:   "Aspirin",
			wantResult: false,
			wantMsg:    models.ErrMedicineNotFound.Error(),
		},
		{
			name: "Person not found",
			snapshot: models.Snapshot{
				People: models.PeopleSlice{},
				Medicines: models.MedicinesMap{
					"Aspirin": &models.MedicineCfg{
						Posology: []models.PosologyEntry{
							{OlderThan: 2 * 365 * 24 * time.Hour}, // 2 years
						},
					},
				},
			},
			person:     "John",
			medicine:   "Aspirin",
			wantResult: false,
			wantMsg:    models.ErrPersonNotFound.Error(),
		},
		{
			name: "Never had a dose",
			snapshot: models.Snapshot{
				People: models.PeopleSlice{
					{Name: "John", Birth: time.Now().AddDate(-10, 0, 0)}, // 10 years old
				},
				Medicines: models.MedicinesMap{
					"Aspirin": &models.MedicineCfg{
						Posology: []models.PosologyEntry{
							{OlderThan: 2 * 365 * 24 * time.Hour}, // 2 years
						},
					},
				},
				Doses: models.DosesMap{},
			},
			person:     "John",
			medicine:   "Aspirin",
			wantResult: true,
			wantMsg:    "they never had a dose",
		},
		{
			name: "Last dose too recent",
			snapshot: models.Snapshot{
				People: models.PeopleSlice{
					{Name: "John", Birth: time.Now().AddDate(-10, 0, 0)}, // 10 years old
				},
				Medicines: models.MedicinesMap{
					"Aspirin": &models.MedicineCfg{
						Posology: []models.PosologyEntry{
							{
								OlderThan:    2 * 365 * 24 * time.Hour, // 2 years
								DoseInterval: 24 * time.Hour,           // 1 day interval
							},
						},
					},
				},
				Doses: models.DosesMap{
					"John": {
						"Aspirin": {time.Now().Add(-12 * time.Hour)}, // 12 hours ago
					},
				},
			},
			person:     "John",
			medicine:   "Aspirin",
			wantResult: false,
			wantMsg:    "their last dose is too recent",
		},
		{
			name: "Too many doses recently",
			snapshot: models.Snapshot{
				People: models.PeopleSlice{
					{Name: "John", Birth: time.Now().AddDate(-10, 0, 0)}, // 10 years old
				},
				Medicines: models.MedicinesMap{
					"Aspirin": &models.MedicineCfg{
						Posology: []models.PosologyEntry{
							{OlderThan: 2 * 365 * 24 * time.Hour, MaxDoses: 1, MaxDosesInterval: 48 * time.Hour}, // 2 years, max 1 dose in 2 days
						},
					},
				},
				Doses: models.DosesMap{
					"John": {
						"Aspirin": {time.Now().Add(-24 * time.Hour)}, // 1 day ago
					},
				},
			},
			person:     "John",
			medicine:   "Aspirin",
			wantResult: false,
			wantMsg:    "they had too many doses recently",
		},
		{
			name: "Can take medicine",
			snapshot: models.Snapshot{
				People: models.PeopleSlice{
					{Name: "John", Birth: time.Now().AddDate(-10, 0, 0)}, // 10 years old
				},
				Medicines: models.MedicinesMap{
					"Aspirin": &models.MedicineCfg{
						Posology: []models.PosologyEntry{
							{OlderThan: 2 * 365 * 24 * time.Hour, DoseInterval: 24 * time.Hour, MaxDoses: 1, MaxDosesInterval: 48 * time.Hour}, // 2 years, 1 day interval, max 1 dose in 2 days
						},
					},
				},
				Doses: models.DosesMap{
					"John": {
						"Aspirin": {time.Now().Add(-72 * time.Hour)}, // 3 days ago
					},
				},
			},
			person:     "John",
			medicine:   "Aspirin",
			wantResult: true,
			wantMsg:    "they haven't had a dose in a while",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResult, gotMsg := tt.snapshot.CanTake(tt.person, tt.medicine)
			if gotResult != tt.wantResult || gotMsg != tt.wantMsg {
				t.Errorf("CanTake() = (%v, %v), want (%v, %v)", gotResult, gotMsg, tt.wantResult, tt.wantMsg)
			}
		})
	}
}
