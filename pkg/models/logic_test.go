package models_test

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"

	"github.com/nanassito/medicine/pkg/models"
)

func TestCanTake(t *testing.T) {
	tests := []struct {
		name         string
		snapshot     models.Snapshot
		person       models.Person
		medicine     models.Medicine
		wantResult   bool
		wantMsg      string
		wantPosology models.PosologyEntry
		wantWaitFor  time.Duration
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
			person:       "John",
			medicine:     "Aspirin",
			wantResult:   true,
			wantMsg:      "they never had a dose",
			wantPosology: models.PosologyEntry{OlderThan: 2 * 365 * 24 * time.Hour},
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
								OlderThan:        2 * 365 * 24 * time.Hour, // 2 years
								DoseInterval:     24 * time.Hour,           // 1 day interval
								MaxDosesInterval: 24 * time.Hour,           // 1 day interval
								MaxDoses:         1000,
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
			wantPosology: models.PosologyEntry{
				OlderThan:        2 * 365 * 24 * time.Hour, // 2 years
				DoseInterval:     24 * time.Hour,           // 1 day interval
				MaxDosesInterval: 24 * time.Hour,           // 1 day interval
				MaxDoses:         1000,
			},
			wantWaitFor: 12 * time.Hour,
		},
		{
			name: "Too many doses recently",
			snapshot: models.Snapshot{
				People: models.PeopleSlice{
					{Name: "John", Birth: time.Now().AddDate(-10, 0, 0)}, // 10 years old
				},
				Medicines: models.MedicinesMap{
					"Aspirin": &models.MedicineCfg{
						Posology: []models.PosologyEntry{ // 37 hours interval
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
			person:       "John",
			medicine:     "Aspirin",
			wantResult:   false,
			wantMsg:      "they had too many doses recently",
			wantPosology: models.PosologyEntry{OlderThan: 2 * 365 * 24 * time.Hour, MaxDoses: 1, MaxDosesInterval: 48 * time.Hour},
			wantWaitFor:  24 * time.Hour,
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
							{OlderThan: 4 * 365 * 24 * time.Hour, DoseInterval: 24 * time.Hour, MaxDoses: 1, MaxDosesInterval: 48 * time.Hour}, // 4 years, 1 day interval, max 1 dose in 2 days
						},
					},
				},
				Doses: models.DosesMap{
					"John": {
						"Aspirin": {time.Now().Add(-72 * time.Hour)}, // 3 days ago
					},
				},
			},
			person:       "John",
			medicine:     "Aspirin",
			wantResult:   true,
			wantMsg:      "they haven't had a dose in a while",
			wantPosology: models.PosologyEntry{OlderThan: 4 * 365 * 24 * time.Hour, DoseInterval: 24 * time.Hour, MaxDoses: 1, MaxDosesInterval: 48 * time.Hour},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResult, gotMsg, posology, waitFor := tt.snapshot.CanTake(tt.person, tt.medicine)
			if diffWaitFor := cmp.Diff(float64(waitFor), float64(tt.wantWaitFor), cmpopts.EquateApprox(0.01, 0)); gotResult != tt.wantResult || gotMsg != tt.wantMsg || posology != tt.wantPosology || diffWaitFor != "" {
				t.Errorf("CanTake() = (%v, %v, %v, %v), want (%v, %v, %v, %v)", gotResult, gotMsg, posology, waitFor, tt.wantResult, tt.wantMsg, tt.wantPosology, tt.wantWaitFor)
			}
		})
	}
}
