package main

import (
	"net/http"

	"github.com/cc-andres-portillo/cronjobs-without-lib-api/handlers"
)

func SetupRoutes() {
	http.HandleFunc("/jobs/add", handlers.AddJobHandler)
	http.HandleFunc("/jobs/list", handlers.ListJobsHandler)
	http.HandleFunc("/jobs/remove", handlers.RemoveJobHandler)
}
