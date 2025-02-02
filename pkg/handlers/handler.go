package handlers

import (
	"errors"
	"fmt"

	"github.com/gorilla/mux"
	"github.com/nanassito/medicine/pkg/models"
	"google.golang.org/api/sheets/v4"
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
	fmt.Println(h.getMedicines())
	fmt.Println(h.getPeople())
	return h, nil
}

func (m *MedicineHandler) getPeople() (map[models.Person]models.PersonCfg, error) {
	val, err := m.GSheetSvc.Spreadsheets.Values.Get(docId, "People!A:C").Do()
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve people from document: %v", err)
	}

	people := make(map[models.Person]models.PersonCfg)
	header := val.Values[0]
	for _, row := range val.Values[1:] {
		name := models.Person(row[0].(string))
		var personCfg models.PersonCfg
		err := models.Unmarshall(row, header, &personCfg)
		if err != nil {
			return nil, err
		}
		people[name] = personCfg
	}

	return people, nil
}

func (m *MedicineHandler) getMedicines() (map[models.Medicine]*models.MedicineCfg, error) {
	val, err := m.GSheetSvc.Spreadsheets.Values.Get(docId, "Medicines!A:G").Do()
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve medicines from document: %v", err)
	}

	medicines := make(map[models.Medicine]*models.MedicineCfg)
	header := val.Values[0]
	for _, row := range val.Values[1:] {
		name := models.Medicine(row[0].(string))
		if _, ok := medicines[name]; !ok {
			medicines[name] = &models.MedicineCfg{Posology: make([]*models.PosologyEntry, 0)}
		}
		var posologyEntry models.PosologyEntry
		err := models.Unmarshall(row, header, &posologyEntry)
		if err != nil {
			return nil, err
		}
		medicine := medicines[name]
		medicine.Posology = append(medicine.Posology, &posologyEntry)
	}

	return medicines, nil
}

func (m *MedicineHandler) getMedicine(name models.Medicine) (*models.MedicineCfg, error) {
	medicines, err := m.getMedicines()
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve medicines: %v", err)
	}
	medicine, ok := medicines[name]
	if !ok {
		return nil, fmt.Errorf("medicine %s not found", name)
	}
	return medicine, nil
}

// func getPosology(person models.PersonCfg, medicine models.MedicineCfg) (models.PosologyEntry, error) {
// 	sort.Slice(medicine.Posology, func(i, j int) bool {
// 		// Sorted from older to younger so the first valid entry is the correct one.
// 		return medicine.Posology[i].OlderThan > medicine.Posology[j].OlderThan
// 	})
// 	for _, entry := range medicine.Posology {
// 		if time.Since(person.Birth) >= entry.OlderThan {
// 			return entry, nil
// 		}
// 	}
// 	return PosologyEntry{}, ErrTooYoung
// }sf

// func nextDose(person PersonCfg, medicineName Medicine, medicine MedicineCfg) (canTakeAfter time.Time, qty string, err error) {
// 	posology, err := getPosology(person, medicine)
// 	if err != nil {
// 		return time.Time{}, "", err
// 	}

// 	if _, ok := person.NextDose[medicineName]; !ok {
// 		// Let's imagine we took a dose at Epoch, just to make sure we have an entry.
// 		person.NextDose[medicineName] = time.Time{}
// 	}
// 	if person.NextDose[medicineName].After(time.Now()) {
// 		err = ErrTooSoon
// 	}
// 	return person.NextDose[medicineName], posology.Quantity, err
// }

// func (h *MedicineHandler) take(w http.ResponseWriter, r *http.Request) {
// 	vars := mux.Vars(r)
// 	slog.Info("take", "vars", vars)
// 	medicineName := Medicine(vars["medicine"])
// 	medicine, ok := h.Medicine[medicineName]
// 	if !ok {
// 		http.Error(w, "Medicine not found", http.StatusNotFound)
// 		return
// 	}
// 	personName := Person(vars["person"])
// 	person, ok := h.People[personName]
// 	if !ok {
// 		http.Error(w, "Person not found", http.StatusNotFound)
// 		return
// 	}

// 	posology, err := getPosology(*person, *medicine)
// 	// We still record the take because it can be dangerous and we would rather have more information than less.
// 	person.NextDose[medicineName] = time.Now().Add(posology.MinInterval)

// 	slog.Info("take",
// 		"who", personName,
// 		"person", person,
// 		"what", medicineName,
// 		"medicine", medicine,
// 		"nextDose", person.NextDose[medicineName],
// 		"err", err,
// 	)

// 	var response string
// 	if err == nil {
// 		response = fmt.Sprintf("Recording %s is taking %s of %s", personName, posology.Quantity, medicineName)
// 	} else {
// 		response = fmt.Sprintf("Error: %v", err)
// 	}
// 	w.Write([]byte(response))
// }

// func (h *MedicineHandler) check(w http.ResponseWriter, r *http.Request) {
// 	vars := mux.Vars(r)
// 	slog.Info("check", "vars", vars)
// 	medicineName := Medicine(vars["medicine"])
// 	medicine, ok := h.Medicine[medicineName]
// 	if !ok {
// 		http.Error(w, "Medicine not found", http.StatusNotFound)
// 		return
// 	}
// 	personName := Person(vars["person"])
// 	person, ok := h.People[personName]
// 	if !ok {
// 		http.Error(w, "Person not found", http.StatusNotFound)
// 		return
// 	}

// 	canTakeAfter, qty, err := nextDose(*person, medicineName, *medicine)
// 	slog.Info("check",
// 		"who", personName,
// 		"person", person,
// 		"what", medicineName,
// 		"medicine", medicine,
// 		"canTakeAfter", canTakeAfter,
// 		"qty", qty,
// 		"err", err,
// 	)
// 	response := fmt.Sprintf("Can Take after: %v\nQuantity: %s\nError: %v", canTakeAfter, qty, err)
// 	w.Write([]byte(response))
// }

// func (h *MedicineHandler) medicineSelect(w http.ResponseWriter, r *http.Request) {
// 	vars := mux.Vars(r)
// 	slog.Info("selection", "vars", vars)
// 	medicineName := Medicine(vars["medicine"])
// 	if _, ok := h.Medicine[medicineName]; !ok {
// 		http.Error(w, "Medicine not found", http.StatusNotFound)
// 		return
// 	}

// 	response := fmt.Sprintf("youpi")
// 	w.Write([]byte(response))
// }

func (h *MedicineHandler) Register(r *mux.Router) {
	// r.HandleFunc("/{medicine}/{person}/take", h.take).Methods(http.MethodGet)
	// r.HandleFunc("/{medicine}/{person}", h.check).Methods(http.MethodGet)
	// r.HandleFunc("/{medicine}", h.medicineSelect).Methods(http.MethodGet)
}
