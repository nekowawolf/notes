package module

import (
	"context"
	"time"
	"fmt"
	"github.com/nekowawolf/notes/config"
	"github.com/nekowawolf/notes/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func InsertOneDoc(collection string, doc interface{}) interface{} {
    insertResult, err := config.Database.Collection(collection).InsertOne(context.TODO(), doc)
    if err != nil {
        fmt.Printf("InsertOneDoc error: %v\n", err)
        return nil
    }
    return insertResult.InsertedID
}


func InsertNote(title, content string) (interface{}, error) {
	note := model.Note{
		ID:      primitive.NewObjectID(),
		Title:   title,
		Content: content,
		Date:    time.Now(),
	}

	insertResult, err := config.Database.Collection("notes").InsertOne(context.TODO(), note)
	if err != nil {
		return nil, err
	}

	return insertResult.InsertedID, nil
}