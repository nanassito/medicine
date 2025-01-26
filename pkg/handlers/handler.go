package handlers

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"sort"
	"time"

	"github.com/gorilla/mux"
)

var (
	ErrTooYoung = errors.New("too young to use this medicine at all")
	ErrTooSoon  = errors.New("too soon to take another dose")
)

type Person string

const (
	Aline  Person = "Aline"
	Dorian Person = "Dorian"
	Zaya   Person = "Zaya"
	Azel   Person = "Azel"
)

type Medicine string

const (
	ChildrenIbuprofen   Medicine = "ChildrenIbuprofen"
	InfantAcetaminophen Medicine = "InfantAcetaminophen"
)

type PosologyEntry struct {
	OlderThan time.Duration
	Interval  time.Duration
	Quantity  string
}

type PersonCfg struct {
	Birth    time.Time
	NextDose map[Medicine]time.Time
}

type MedicineCfg struct {
	Posology []PosologyEntry
}

type MedicineHandler struct {
	People   map[Person]*PersonCfg
	Medicine map[Medicine]*MedicineCfg
}

func getPosology(person PersonCfg, medicine MedicineCfg) (PosologyEntry, error) {
	sort.Slice(medicine.Posology, func(i, j int) bool {
		// Sorted from older to younger so the first valid entry is the correct one.
		return medicine.Posology[i].OlderThan > medicine.Posology[j].OlderThan
	})
	for _, entry := range medicine.Posology {
		if time.Since(person.Birth) >= entry.OlderThan {
			return entry, nil
		}
	}
	return PosologyEntry{}, ErrTooYoung
}

func nextDose(person PersonCfg, medicineName Medicine, medicine MedicineCfg) (canTakeAfter time.Time, qty string, err error) {
	posology, err := getPosology(person, medicine)
	if err != nil {
		return time.Time{}, "", err
	}

	if _, ok := person.NextDose[medicineName]; !ok {
		// Let's imagine we took a dose at Epoch, just to make sure we have an entry.
		person.NextDose[medicineName] = time.Time{}
	}
	if person.NextDose[medicineName].After(time.Now()) {
		err = ErrTooSoon
	}
	return person.NextDose[medicineName], posology.Quantity, err
}

func (h *MedicineHandler) take(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slog.Info("take", "vars", vars)
	medicineName := Medicine(vars["medicine"])
	medicine, ok := h.Medicine[medicineName]
	if !ok {
		http.Error(w, "Medicine not found", http.StatusNotFound)
		return
	}
	personName := Person(vars["person"])
	person, ok := h.People[personName]
	if !ok {
		http.Error(w, "Person not found", http.StatusNotFound)
		return
	}

	posology, err := getPosology(*person, *medicine)
	// We still record the take because it can be dangerous and we would rather have more information than less.
	person.NextDose[medicineName] = time.Now().Add(posology.Interval)

	slog.Info("take",
		"who", personName,
		"person", person,
		"what", medicineName,
		"medicine", medicine,
		"nextDose", person.NextDose[medicineName],
		"err", err,
	)

	var response string
	if err == nil {
		response = fmt.Sprintf("Recording %s is taking %s of %s", personName, posology.Quantity, medicineName)
	} else {
		response = fmt.Sprintf("Error: %v", err)
	}
	w.Write([]byte(response))
}

func (h *MedicineHandler) check(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slog.Info("check", "vars", vars)
	medicineName := Medicine(vars["medicine"])
	medicine, ok := h.Medicine[medicineName]
	if !ok {
		http.Error(w, "Medicine not found", http.StatusNotFound)
		return
	}
	personName := Person(vars["person"])
	person, ok := h.People[personName]
	if !ok {
		http.Error(w, "Person not found", http.StatusNotFound)
		return
	}

	canTakeAfter, qty, err := nextDose(*person, medicineName, *medicine)
	slog.Info("check",
		"who", personName,
		"person", person,
		"what", medicineName,
		"medicine", medicine,
		"canTakeAfter", canTakeAfter,
		"qty", qty,
		"err", err,
	)
	response := fmt.Sprintf("Can Take after: %v\nQuantity: %s\nError: %v", canTakeAfter, qty, err)
	w.Write([]byte(response))
}

func (h *MedicineHandler) medicineSelect(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slog.Info("selection", "vars", vars)
	medicineName := Medicine(vars["medicine"])
	if _, ok := h.Medicine[medicineName]; !ok {
		http.Error(w, "Medicine not found", http.StatusNotFound)
		return
	}

	response := fmt.Sprintf("youpi")
	w.Write([]byte(response))
}

func (h *MedicineHandler) Register(r *mux.Router) {
	r.HandleFunc("/{medicine}/{person}/take", h.take).Methods(http.MethodGet)
	r.HandleFunc("/{medicine}/{person}", h.check).Methods(http.MethodGet)
	r.HandleFunc("/{medicine}", h.medicineSelect).Methods(http.MethodGet)
}
