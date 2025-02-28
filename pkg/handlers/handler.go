package handlers

import (
	"fmt"
	"log/slog"
	"net/http"
	"sort"
	"time"

	"github.com/gorilla/mux"

	"github.com/nanassito/medicine/pkg/models"
	"github.com/nanassito/medicine/pkg/templates"
)

func (h *MedicineHandler) medicineOverview(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slog.Info("selection", "vars", vars)

	snapshot, err := h.getAll(r.Context())
	if err != nil {
		http.Error(w, fmt.Sprintf("unable to retrieve data: %v", err), http.StatusInternalServerError)
		return
	}
	medicineName := models.Medicine(vars["medicine"])
	if !snapshot.HasMedicine(medicineName) {
		http.Error(w, fmt.Sprintf("medicine %s not found", medicineName), http.StatusNotFound)
		return
	}

	data := struct {
		MedicineName models.Medicine
		People       []models.PersonCfg
	}{
		MedicineName: medicineName,
		People:       snapshot.People,
	}
	if err = templates.MedicineOverview.Execute(w, data); err != nil {
		http.Error(w, fmt.Sprintf("unable to execute template: %v", err), http.StatusInternalServerError)
	}
}

func (h *MedicineHandler) take(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slog.Info("selection", "vars", vars)

	snapshot, err := h.getAll(r.Context())
	if err != nil {
		http.Error(w, fmt.Sprintf("unable to retrieve data: %v", err), http.StatusInternalServerError)
		return
	}
	medicineName := models.Medicine(vars["medicine"])
	if !snapshot.HasMedicine(medicineName) {
		http.Error(w, fmt.Sprintf("medicine %s not found", medicineName), http.StatusNotFound)
		return
	}

	personName := models.Person(vars["person"])
	if !snapshot.HasPerson(personName) {
		http.Error(w, fmt.Sprintf("person %s not found", personName), http.StatusNotFound)
		return
	}

	if err = h.logDoseIntake(personName, medicineName); err != nil {
		http.Error(w, fmt.Sprintf("unable to register that %s was taken by %s: %v", medicineName, personName, err), http.StatusInternalServerError)
		return
	}

	canTake, reason, posology, waitFor := snapshot.CanTake(personName, medicineName)
	if !canTake {
		if waitFor > time.Duration(0.9*float64(posology.DoseInterval)) {
			w.Write([]byte(fmt.Sprintf("Do NOT take this ! %s", reason)))
			return
		}
	}

	http.Redirect(w, r, fmt.Sprintf("/%s", medicineName), http.StatusSeeOther)
}

func (h *MedicineHandler) list(w http.ResponseWriter, r *http.Request) {
	snapshot, err := h.getAll(r.Context())
	if err != nil {
		http.Error(w, fmt.Sprintf("unable to retrieve data: %v", err), http.StatusInternalServerError)
		return
	}

	medicines := make([]string, 0)
	for medicine := range snapshot.Medicines {
		medicines = append(medicines, string(medicine))
	}
	sort.Strings(medicines)
	data := struct {
		Medicines []string
	}{
		Medicines: medicines,
	}
	if err = templates.List.Execute(w, data); err != nil {
		http.Error(w, fmt.Sprintf("unable to execute template: %v", err), http.StatusInternalServerError)
	}
}

func (h *MedicineHandler) medicineFor(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slog.Info("selection", "vars", vars)

	snapshot, err := h.getAll(r.Context())
	if err != nil {
		http.Error(w, fmt.Sprintf("unable to retrieve data: %v", err), http.StatusInternalServerError)
		return
	}
	medicineName := models.Medicine(vars["medicine"])
	if !snapshot.HasMedicine(medicineName) {
		http.Error(w, fmt.Sprintf("medicine %s not found", medicineName), http.StatusNotFound)
		return
	}

	personName := models.Person(vars["person"])
	if !snapshot.HasPerson(personName) {
		http.Error(w, fmt.Sprintf("person %s not found", personName), http.StatusNotFound)
		return
	}

	canTake, reason, posology, waitFor := snapshot.CanTake(personName, medicineName)

	data := struct {
		MedicineName models.Medicine
		Who          models.PersonCfg
		Reason       string
		CanTake      bool
		Posology     models.PosologyEntry
		WaitForPct   float64
		WaitFor      time.Duration
	}{
		MedicineName: medicineName,
		Who:          snapshot.GetPerson(personName),
		Reason:       reason,
		CanTake:      canTake,
		Posology:     posology,
		WaitForPct:   float64(waitFor) / float64(posology.DoseInterval),
		WaitFor:      waitFor,
	}
	if err = templates.MedicineFor.Execute(w, data); err != nil {
		http.Error(w, fmt.Sprintf("unable to execute template: %v", err), http.StatusInternalServerError)
	}
}

func (h *MedicineHandler) Register(r *mux.Router) {
	r.HandleFunc("/{medicine}/{person}/take", h.take).Methods(http.MethodGet)
	r.HandleFunc("/{medicine}/{person}", h.medicineFor).Methods(http.MethodGet)
	r.HandleFunc("/{medicine}", h.medicineOverview).Methods(http.MethodGet)
	r.HandleFunc("/", h.list).Methods(http.MethodGet)
}
