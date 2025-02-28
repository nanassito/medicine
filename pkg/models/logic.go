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

func (s *Snapshot) GetPerson(person Person) PersonCfg {
	for _, p := range s.People {
		if p.Name == person {
			return p
		}
	}
	return PersonCfg{}
}

func (s *Snapshot) CanTake(who Person, what Medicine) (canTake bool, reason string, posology PosologyEntry, waitFor time.Duration) {
	posology, err := s.GetPosology(who, what)
	waitFor = time.Duration(0)
	reason = "This is a bug"
	canTake = false

	if err != nil {
		reason = err.Error()
		// If the person is too young or there is some missing data we'll get an error.
		return
	}

	// This person never had a dose so it's fine.
	if _, ok := s.Doses[who][what]; !ok {
		canTake = true
		reason = "they never had a dose"
		return
	}

	doses := []time.Time{}
	sort.Slice(s.Doses[who][what], func(i, j int) bool {
		// Sort from most recent to oldest
		return s.Doses[who][what][i].After(s.Doses[who][what][j])
	})
	for _, dose := range s.Doses[who][what] {
		if dose.After(time.Now().Add(-posology.MaxDosesInterval)) {
			doses = append(doses, dose)
		}
	}
	// Ok now we know the person can generally take this medicine.
	// Let's check if they haven't over done it.
	canTake = true
	reason = "they haven't had a dose in a while"

	if len(doses) > 0 && time.Since(doses[0]) <= posology.DoseInterval {
		canTake = false
		reason = "their last dose is too recent"
		waitFor = posology.DoseInterval - time.Since(doses[0])
	}

	if len(doses) > 0 && len(doses) >= int(posology.MaxDoses) {
		canTake = false
		reason = "they had too many doses recently"
		oldestRelevantDose := doses[int(posology.MaxDoses)-1]
		waitFor = max(waitFor, posology.MaxDosesInterval-time.Since(oldestRelevantDose))
	}

	return
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
