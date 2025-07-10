package job

import (
	"context"
	"fmt"
	"time"

	"github.com/cc-andres-portillo/cronjobs-without-lib-api/db"
	"github.com/cc-andres-portillo/cronjobs-without-lib-api/models"
	"go.mongodb.org/mongo-driver/bson"
)

func NotifyHotLeads() {
	fmt.Println("🚨 notifyHotLeads running:", time.Now())

	// Buscar leads vencidos
	now := time.Now()
	filter := bson.M{
		"isRemove":       false,
		"statusId":       bson.M{"$in": []string{"status_contacting_id"}},
		"time.dateLimit": bson.M{"$lte": now},
	}

	cursor, err := db.CronjobsCollection.Database().Collection("leads").Find(context.TODO(), filter)
	if err != nil {
		fmt.Println("❌ Error finding leads:", err)
		return
	}
	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		var lead models.Lead
		if err := cursor.Decode(&lead); err != nil {
			continue
		}

		// Actualizar lead
		update := bson.M{
			"$set": bson.M{
				"hot":                           true,
				"statusId":                      "status_open_id",
				"available":                     true,
				"settings.contacting.moveCount": 0,
				"settings.interest.isLocked":    false,
				"settings.interest.lockedBy":    "",
				"userId":                        "",
				"users":                         []string{},
				"selected":                      false,
				"producerUserId":                "",
				"financerUserId":                "",
			},
			"$unset": bson.M{
				"agencyId":         "",
				"producerTeamId":   "",
				"teamId":           "",
				"leaderUserId":     "",
				"supervisorUserId": "",
			},
			"$push": bson.M{
				"previousProducerIds": lead.ProducerUserID,
			},
		}

		_, err := db.CronjobsCollection.Database().Collection("leads").UpdateByID(context.TODO(), lead.ID, update)
		if err != nil {
			fmt.Println("❌ Error updating lead:", lead.ID, err)
		}

		// Agregar logs, historia y notificaciones si se requiere
	}
}
