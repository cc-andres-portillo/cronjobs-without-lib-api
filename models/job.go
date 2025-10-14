package models

type Job struct {
	ID       string `json:"id" bson:"_id"`
	Message  string `json:"message" bson:"message"`
	Interval int    `json:"interval" bson:"interval"` // en segundos
}
