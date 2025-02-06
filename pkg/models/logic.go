package models

import (
	"errors"
	"sort"
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

func (s *Snapshot) HasPerson(person Person) bool {
	for _, p := range s.People {
		if p.Name == person {
			return true
		}
	}
	return false
}

func (s *Snapshot) CanTake(who Person, what Medicine) (bool, string, PosologyEntry) {
	posology, err := s.GetPosology(who, what)
	if err != nil {
		// If the person is too young or there is some missing data we'll get an error.
		return false, err.Error(), posology
	}

	// This person never had a dose so it's fine.
	if _, ok := s.Doses[who][what]; !ok {
		return true, "they never had a dose", posology
	}

	var lastDose time.Time
	var numDoses int64
	sort.Slice(s.Doses[who][what], func(i, j int) bool {
		return s.Doses[who][what][i].After(s.Doses[who][what][j])
	})
	for _, dose := range s.Doses[who][what] {
		if dose.After(lastDose) {
			lastDose = dose
		}
		if dose.After(time.Now().Add(-posology.MaxDosesInterval)) {
			numDoses++
		}
	}
	if time.Since(lastDose) <= posology.DoseInterval {
		return false, "their last dose is too recent", posology
	}
	if numDoses >= posology.MaxDoses {
		return false, "they had too many doses recently", posology
	}

	return true, "they haven't had a dose in a while", posology
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

	sort.Slice(medicine.Posology, func(i, j int) bool {
		if medicine.Posology[i].OlderThan == medicine.Posology[j].OlderThan {
			return medicine.Posology[i].HeavierThan > medicine.Posology[j].HeavierThan
		}
		return medicine.Posology[i].OlderThan > medicine.Posology[j].OlderThan
	})

	for _, entry := range medicine.Posology {
		if time.Since(person.Birth) >= entry.OlderThan || (person.Weight >= entry.HeavierThan && entry.HeavierThan > 0) {
			return entry, nil
		}
	}

	return PosologyEntry{}, ErrTooYoung
}
