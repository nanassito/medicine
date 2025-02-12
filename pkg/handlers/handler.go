package handlers

import (
	"fmt"
	"log/slog"
	"net/http"
	"sort"

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
		CanTake      []struct {
			Who     models.PersonCfg
			CanTake bool
			Reason  string
			Dose    string
		}
	}{
		MedicineName: medicineName,
		CanTake: make([]struct {
			Who     models.PersonCfg
			CanTake bool
			Reason  string
			Dose    string
		}, 0),
	}

	for _, person := range snapshot.People {
		canTake, reason, posology := snapshot.CanTake(person.Name, medicineName)
		data.CanTake = append(data.CanTake, struct {
			Who     models.PersonCfg
			CanTake bool
			Reason  string
			Dose    string
		}{person, canTake, reason, posology.Dose})
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

	canTake, reason, _ := snapshot.CanTake(personName, medicineName)
	if !canTake {
		w.Write([]byte(fmt.Sprintf("Do NOT take this ! %s", reason)))
		return
	}

	if err = h.logDoseIntake(personName, medicineName); err != nil {
		http.Error(w, fmt.Sprintf("unable to register that %s was taken by %s: %v", medicineName, personName, err), http.StatusInternalServerError)
		return
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

func (h *MedicineHandler) Register(r *mux.Router) {
	r.HandleFunc("/{medicine}/{person}", h.take).Methods(http.MethodGet)
	r.HandleFunc("/{medicine}", h.medicineOverview).Methods(http.MethodGet)
	r.HandleFunc("/", h.list).Methods(http.MethodGet)
}
