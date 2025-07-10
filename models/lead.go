package models

import "time"

type Lead struct {
	ID             string    `bson:"_id"`
	Name           string    `bson:"name"`
	Company        string    `bson:"company"`
	USDOT          string    `bson:"usdot"`
	StatusID       string    `bson:"statusId"`
	AgencyID       string    `bson:"agencyId"`
	UserID         string    `bson:"userId"`
	TeamID         string    `bson:"teamId"`
	ProducerUserID string    `bson:"producerUserId"`
	ProducerTeamID string    `bson:"producerTeamId"`
	DateLimit      time.Time `bson:"time.dateLimit"`
}
