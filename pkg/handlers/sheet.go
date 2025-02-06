package handlers

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sort"
	"time"

	"golang.org/x/sync/errgroup"
	"google.golang.org/api/sheets/v4"

	"github.com/nanassito/medicine/pkg/models"
)

const (
	docId = "1MGRP9e0aUBvukLeo2oAP1WP4NM7mBkGn1Z0wyakOdbo"
)

var (
	ErrTooYoung = errors.New("too young to use this medicine at all")
	ErrTooSoon  = errors.New("too soon to take another dose")
)

type MedicineHandler struct {
	GSheetSvc *sheets.Service
}

func NewMedicineHandler(svc *sheets.Service) (*MedicineHandler, error) {
	_, err := svc.Spreadsheets.Get(docId).Do()
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve document: %v", err)
	}
	h := &MedicineHandler{GSheetSvc: svc}
	return h, nil
}

func (m *MedicineHandler) getPeople() (models.PeopleSlice, error) {
	val, err := m.GSheetSvc.Spreadsheets.Values.Get(docId, "People!A:D").Do()
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve people from document: %v", err)
	}

	people := make(models.PeopleSlice, 0)
	header := val.Values[0]
	for _, row := range val.Values[1:] {
		var personCfg models.PersonCfg
		err := models.Unmarshall(row, header, &personCfg)
		if err != nil {
			return nil, err
		}
		people = append(people, personCfg)
	}

	return people, nil
}

func (m *MedicineHandler) getDoses() (models.DosesMap, error) {
	val, err := m.GSheetSvc.Spreadsheets.Values.Get(docId, "Events!A:C").Do()
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve doses from document: %v", err)
	}

	doses := make(models.DosesMap)
	header := val.Values[0]
	for _, row := range val.Values[1:] {
		var dose models.Dose
		err := models.Unmarshall(row, header, &dose)
		if err != nil {
			return nil, err
		}
		if _, ok := doses[dose.Who]; !ok {
			doses[dose.Who] = make(map[models.Medicine][]time.Time)
		}
		if _, ok := doses[dose.Who][dose.What]; !ok {
			doses[dose.Who][dose.What] = make([]time.Time, 0)
		}
		doses[dose.Who][dose.What] = append(doses[dose.Who][dose.What], dose.When)
	}

	for _, personDoses := range doses {
		for _, medicineDoses := range personDoses {
			sort.Slice(medicineDoses, func(i, j int) bool {
				return medicineDoses[i].After(medicineDoses[j])
			})
		}
	}

	return doses, nil
}

func (m *MedicineHandler) getMedicines() (models.MedicinesMap, error) {
	val, err := m.GSheetSvc.Spreadsheets.Values.Get(docId, "Medicines!A:G").Do()
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve medicines from document: %v", err)
	}

	medicines := make(models.MedicinesMap)
	header := val.Values[0]
	for _, row := range val.Values[1:] {
		name := models.Medicine(row[0].(string))
		if _, ok := medicines[name]; !ok {
			medicines[name] = &models.MedicineCfg{Posology: make([]models.PosologyEntry, 0)}
		}
		var posologyEntry models.PosologyEntry
		err := models.Unmarshall(row, header, &posologyEntry)
		if err != nil {
			return nil, err
		}
		medicine := medicines[name]
		medicine.Posology = append(medicine.Posology, posologyEntry)
	}

	for _, medicine := range medicines {
		// Reverse sort so we can iterate from older to younger.
		sort.Slice(medicine.Posology, func(i, j int) bool {
			if medicine.Posology[i].OlderThan == medicine.Posology[j].OlderThan {
				return medicine.Posology[i].HeavierThan > medicine.Posology[j].HeavierThan
			}
			return medicine.Posology[i].OlderThan > medicine.Posology[j].OlderThan
		})
	}

	return medicines, nil
}

func (m *MedicineHandler) getAll(ctx context.Context) (snapshot models.Snapshot, err error) {
	group, _ := errgroup.WithContext(ctx)
	group.Go(func() error {
		people, err := m.getPeople()
		if err != nil {
			return fmt.Errorf("unable to retrieve people: %v", err)
		}
		snapshot.People = people
		return nil
	})
	group.Go(func() error {
		medicines, err := m.getMedicines()
		if err != nil {
			return fmt.Errorf("unable to retrieve medicines: %v", err)
		}
		snapshot.Medicines = medicines
		return nil
	})
	group.Go(func() error {
		doses, err := m.getDoses()
		if err != nil {
			return fmt.Errorf("unable to retrieve doses: %v", err)
		}
		snapshot.Doses = doses
		return nil
	})
	if err := group.Wait(); err != nil {
		return snapshot, err
	}

	return snapshot, nil
}

func (m *MedicineHandler) logDoseIntake(personName models.Person, medicineName models.Medicine) error {
	slog.Info("dose intake", "person", personName, "medicine", medicineName)
	resp, err := m.GSheetSvc.Spreadsheets.Values.Append(docId, "Events!A2", &sheets.ValueRange{
		Values: [][]interface{}{
			{personName, medicineName, time.Now().UTC().Format(time.DateTime)},
		},
	}).InsertDataOption("INSERT_ROWS").ValueInputOption("USER_ENTERED").Do()
	if err != nil {
		slog.Error("unable to log dose intake", "error", err)
		return fmt.Errorf("unable to log dose intake: %v", err)
	}
	fmt.Println(resp)
	return nil
}
