package models

type Job struct {
	ID       string `json:"id"`
	Message  string `json:"message"`
	Interval int    `json:"interval"` // En segundos
}
