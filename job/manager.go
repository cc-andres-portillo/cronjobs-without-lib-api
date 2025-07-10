package job

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/cc-andres-portillo/cronjobs-without-lib-api/db"
	"github.com/cc-andres-portillo/cronjobs-without-lib-api/models"
	"go.mongodb.org/mongo-driver/bson"
)

type JobManager struct {
	mu   sync.Mutex
	jobs map[string]chan bool
}

var Manager = JobManager{
	jobs: make(map[string]chan bool),
}

func (jm *JobManager) AddJob(job models.Job) {
	jm.mu.Lock()
	defer jm.mu.Unlock()

	if _, exists := jm.jobs[job.ID]; exists {
		log.Printf("[WARN] Job %s already exists", job.ID)
		return
	}

	stopChan := make(chan bool)
	jm.jobs[job.ID] = stopChan

	// Guardar en MongoDB
	_, err := db.CronjobsCollection.InsertOne(context.TODO(), job)
	if err != nil {
		log.Printf("❌ Failed to save job %s: %v", job.ID, err)
	}

	// Obtener función del registro si está definida
	execFunc, found := JobRegistry[job.ID]
	if !found {
		// fallback genérico si no hay función registrada
		execFunc = func() {
			fmt.Printf("[Job: %s] %s\n", job.ID, job.Message)
		}
	}

	// Ejecutar job usando ticker
	go func() {
		ticker := time.NewTicker(time.Duration(job.Interval) * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				execFunc()
			case <-stopChan:
				fmt.Printf("[Job: %s] Stopped\n", job.ID)
				return
			}
		}
	}()

}

func (jm *JobManager) RemoveJob(id string) {
	jm.mu.Lock()
	defer jm.mu.Unlock()

	if stopChan, exists := jm.jobs[id]; exists {
		stopChan <- true
		delete(jm.jobs, id)

		_, err := db.CronjobsCollection.DeleteOne(context.TODO(), bson.M{"_id": id})
		if err != nil {
			log.Printf("❌ Failed to delete job %s from DB: %v", id, err)
		}
	}
}

func (jm *JobManager) ListJobs() []string {
	jm.mu.Lock()
	defer jm.mu.Unlock()

	ids := []string{}
	for id := range jm.jobs {
		ids = append(ids, id)
	}
	return ids
}

func (jm *JobManager) LoadJobsFromDB() {
	cursor, err := db.CronjobsCollection.Find(context.TODO(), bson.M{})
	if err != nil {
		log.Println("❌ Error loading jobs:", err)
		return
	}
	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		var job models.Job
		if err := cursor.Decode(&job); err == nil {
			jm.AddJob(job)
		}
	}
	log.Println("✅ Jobs loaded from MongoDB")
}

func (jm *JobManager) AutoRegisterRegistryJobs() {
	for jobID := range JobRegistry {
		// Verificar si ya existe en MongoDB
		count, err := db.CronjobsCollection.CountDocuments(context.TODO(), bson.M{"_id": jobID})
		if err != nil {
			fmt.Println("❌ Error checking job:", jobID, err)
			continue
		}
		if count == 0 {
			// Agregar job a MongoDB y al manager
			newJob := models.Job{
				ID:       jobID,
				Message:  jobID, // puedes usar otro mensaje si quieres
				Interval: 300,   // default 5 min (o cambiar por tabla/mapa si quieres)
			}
			jm.AddJob(newJob)
			fmt.Println("✅ Auto-registered job:", jobID)
		} else {
			fmt.Println("ℹ️ Job already exists:", jobID)
		}
	}
}
