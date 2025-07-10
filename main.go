package main

import (
	"fmt"
	"net/http"

	"github.com/cc-andres-portillo/cronjobs-without-lib-api/db"
	"github.com/cc-andres-portillo/cronjobs-without-lib-api/job"
)

func main() {
	fmt.Println("🚀 Starting CronJob API on port 8080...")

	db.InitMongo()

	// Auto-registrar jobs definidos en el código
	job.Manager.AutoRegisterRegistryJobs()

	job.Manager.LoadJobsFromDB()

	SetupRoutes()
	http.ListenAndServe(":8080", nil)
}
