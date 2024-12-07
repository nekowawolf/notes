package module

import (
	"context"
	"time"
	"fmt"
	"github.com/nekowawolf/notes/config"
	"github.com/nekowawolf/notes/model"
	"github.com/nekowawolf/notes/auth"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/bson"
)

func InsertOneDoc(collection string, doc interface{}) interface{} {
    insertResult, err := config.Database.Collection(collection).InsertOne(context.TODO(), doc)
    if err != nil {
        fmt.Printf("InsertOneDoc error: %v\n", err)
        return nil
    }
    return insertResult.InsertedID
}

func InsertNotes(title, content string) (interface{}, error) {
	notes := model.Notes{
		ID:      primitive.NewObjectID(),
		Title:   title,
		Content: content,
		Date:    time.Now(),
	}

	insertResult, err := config.Database.Collection("notes").InsertOne(context.TODO(), notes)
	if err != nil {
		return nil, err
	}

	return insertResult.InsertedID, nil
}

func InsertAdmin(username, password string) (interface{}, error) {
    hashedPassword, err := auth.HashPassword(password)
    if err != nil {
        return nil, err
    }

    admin := model.Admin{
        ID:       primitive.NewObjectID(),
        Username: username,
        Password: hashedPassword,
    }

    insertResult, err := config.Database.Collection("admin").InsertOne(context.TODO(), admin)
    if err != nil {
        return nil, err
    }

    return insertResult.InsertedID, nil
}

func LoginAdmin(username, password string) (model.Admin, error) {
    var admin model.Admin

    err := config.Database.Collection("admin").FindOne(context.TODO(), bson.M{"username": username}).Decode(&admin)
    if err != nil {
        return model.Admin{}, fmt.Errorf("admin not found: %v", err)
    }

    // Compare password
    if err := auth.VerifyPassword(admin.Password, password); err != nil {
        return model.Admin{}, fmt.Errorf("invalid password: %v", err)
    }

    return admin, nil
}


func GetAllNotes() ([]model.Notes, error) {
	collection := config.Database.Collection("notes")

	cursor, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		return nil, fmt.Errorf("GetAllNotes Find: %v", err)
	}

	var notes []model.Notes
	if err = cursor.All(context.TODO(), &notes); err != nil {
		return nil, fmt.Errorf("GetAllNotes All: %v", err)
	}

	return notes, nil
}

func GetNotesByID(id primitive.ObjectID) (model.Notes, error) {
	collection := config.Database.Collection("notes")
	var notes model.Notes
	err := collection.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&notes)
	if err != nil {
		return model.Notes{}, err
	}
	return notes, nil
}

func UpdateNotesByID(id primitive.ObjectID, title string, content string) (model.Notes, error) {
	collection := config.Database.Collection("notes")

	update := bson.M{
		"$set": bson.M{
			"title":   title,
			"content": content,
		},
	}

	filter := bson.M{"_id": id}

	var updatedNote model.Notes
	err := collection.FindOneAndUpdate(context.TODO(), filter, update).Decode(&updatedNote)
	if err != nil {
		return model.Notes{}, err
	}

	return updatedNote, nil
}

func DeleteNotesByID(id primitive.ObjectID) error {
    collection := config.Database.Collection("notes")
    filter := bson.M{"_id": id}

    result, err := collection.DeleteOne(context.TODO(), filter)
    if err != nil {
        return fmt.Errorf("error deleting note for ID %s: %s", id.Hex(), err.Error())
    }

    if result.DeletedCount == 0 {
        return fmt.Errorf("no note found with ID %s", id.Hex())
    }

    return nil
}



