package models

import (
	"errors"
	"time"
)

var (
	ErrMedicineNotFound = errors.New("<error: medicine not found>")
	ErrPersonNotFound   = errors.New("<error: person not found>")
	ErrTooYoung         = errors.New("they are too young")
)

type PeopleSlice []PersonCfg

type DosesMap map[Person]map[Medicine][]time.Time

type MedicinesMap map[Medicine]*MedicineCfg

type Snapshot struct {
	People    PeopleSlice
	Doses     DosesMap
	Medicines MedicinesMap
}

func (s *Snapshot) HasMedicine(medicine Medicine) bool {
	_, ok := s.Medicines[medicine]
	return ok
}

func (s *Snapshot) CanTake(who Person, what Medicine) (bool, string) {
	posology, err := s.GetPosology(who, what)
	if err != nil {
		// If the person is too young or there is some missing data we'll get an error.
		return false, err.Error()
	}

	// This person never had a dose so it's fine.
	if _, ok := s.Doses[who][what]; !ok {
		return true, "they never had a dose"
	}

	var lastDose time.Time
	var numDoses int64
	for _, dose := range s.Doses[who][what] {
		if dose.After(lastDose) {
			lastDose = dose
		}
		if dose.After(time.Now().Add(-posology.MaxDosesInterval)) {
			numDoses++
		}
	}
	if time.Since(lastDose) < posology.DoseInterval {
		return false, "their last dose is too recent"
	}
	if numDoses >= posology.MaxDoses {
		return false, "they had too many doses recently"
	}

	return true, "they haven't had a dose in a while"
}

func (s *Snapshot) GetPosology(personName Person, medicineName Medicine) (PosologyEntry, error) {
	medicine, ok := s.Medicines[medicineName]
	if !ok {
		return PosologyEntry{}, ErrMedicineNotFound
	}
	var person *PersonCfg
	for _, candidate := range s.People {
		if candidate.Name == personName {
			person = &candidate
		}
	}
	if person == nil {
		return PosologyEntry{}, ErrPersonNotFound
	}

	for _, entry := range medicine.Posology {
		if time.Since(person.Birth) >= entry.OlderThan || person.Weight >= entry.HeavierThan {
			return entry, nil
		}
	}

	return PosologyEntry{}, ErrTooYoung
}
