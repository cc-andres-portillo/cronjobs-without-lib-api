package handlers

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/cc-andres-portillo/cronjobs-without-lib-api/db"
	"github.com/cc-andres-portillo/cronjobs-without-lib-api/job"
	"github.com/cc-andres-portillo/cronjobs-without-lib-api/models"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
)

func AddJobHandler(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)
	var input struct {
		Message  string `json:"message"`
		Interval int    `json:"interval"`
	}

	if err := json.Unmarshal(body, &input); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Verificar si ya existe un job con el mismo mensaje + intervalo
	filter := bson.M{
		"message":  input.Message,
		"interval": input.Interval,
	}

	count, err := db.CronjobsCollection.CountDocuments(context.TODO(), filter)
	if err != nil {
		http.Error(w, "Error checking existing jobs", http.StatusInternalServerError)
		return
	}

	if count > 0 {
		http.Error(w, "Job already exists with same message and interval", http.StatusConflict)
		return
	}

	// Crear nuevo job
	newJob := models.Job{
		ID:       uuid.New().String(),
		Message:  input.Message,
		Interval: input.Interval,
	}

	job.Manager.AddJob(newJob)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newJob)
}

func ListJobsHandler(w http.ResponseWriter, r *http.Request) {
	jobIDs := job.Manager.ListJobs()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(jobIDs)
}

func RemoveJobHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	id := query.Get("id")

	if id == "" {
		http.Error(w, "Missing job ID", http.StatusBadRequest)
		return
	}

	job.Manager.RemoveJob(id)
	w.WriteHeader(http.StatusNoContent)
}
