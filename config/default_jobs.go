package config

import "github.com/cc-andres-portillo/cronjobs-without-lib-api/models"

func DefaultJobs() []models.Job {
	return []models.Job{
		{
			ID:       "startup-job-1",
			Message:  "Este es un cronjob automático",
			Interval: 5,
		},
		{
			ID:       "startup-job-2",
			Message:  "Otro cronjob precargado",
			Interval: 10,
		},
	}
}
