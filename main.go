package main

import (
	"fmt"
	"net/http"

	"github.com/cc-andres-portillo/cronjobs-without-lib-api/config"
	"github.com/cc-andres-portillo/cronjobs-without-lib-api/job"
)

func main() {
	fmt.Println("Starting CronJob API on port 8080...")

	// Ejecutar jobs por defecto
	for _, j := range config.DefaultJobs() {
		job.Manager.AddJob(j)
	}

	// Configurar las rutas y levantar el servidor
	SetupRoutes()
	http.ListenAndServe(":8080", nil)
}
