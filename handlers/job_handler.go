package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/cc-andres-portillo/cronjobs-without-lib-api/job"
	"github.com/cc-andres-portillo/cronjobs-without-lib-api/models"

	"github.com/google/uuid"
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
