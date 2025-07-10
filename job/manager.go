package job

import (
	"fmt"
	"sync"
	"time"

	"github.com/cc-andres-portillo/cronjobs-without-lib-api/models"
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

	stopChan := make(chan bool)
	jm.jobs[job.ID] = stopChan

	go func() {
		ticker := time.NewTicker(time.Duration(job.Interval) * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				fmt.Printf("[Job: %s] %s\n", job.ID, job.Message)
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
