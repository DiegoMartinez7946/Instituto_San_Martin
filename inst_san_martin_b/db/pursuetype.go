package db

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/benjacifre10/san_martin_b/config"
	"github.com/benjacifre10/san_martin_b/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

/***************************************************************/
/***************************************************************/
/* GetPursueTypesDB get the pursue types from db */
func GetPursueTypesDB() ([]*models.PursueType, bool) {
	ctx, cancel := context.WithTimeout(context.Background(), 15 * time.Second)
	defer cancel()

	db := config.MongoConnection.Database("san_martin")
	collection := db.Collection("pursue_type")

	var results []*models.PursueType

	condition := bson.M {  }
	optionsQuery := options.Find()
	optionsQuery.SetSort(bson.D {{ Key: "type", Value: -1}})

	pursueTypes, err := collection.Find(ctx, condition, optionsQuery)
	if err != nil {
		log.Fatal(err.Error())
		return results, false
	}

	for pursueTypes.Next(context.TODO()) {
		var row models.PursueType
		err := pursueTypes.Decode(&row)
		if err != nil {
			return results, false
		}
		results = append(results, &row)
	}

	return results, true
}

/***************************************************************/
/***************************************************************/
/* InsertPursueTypeDB insert one pursue type in db */
func InsertPursueTypeDB(p models.PursueType) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15 * time.Second)
	defer cancel()

	db := config.MongoConnection.Database("san_martin")
	collection := db.Collection("pursue_type")

	row := bson.M {
		"type": p.Type,
	}

	result, err := collection.InsertOne(ctx, row)
	if err != nil {
		return "Hubo un error al insertar la modalidad de cursado", err
	}
	
	objID, _ := result.InsertedID.(primitive.ObjectID)
	return objID.Hex(), nil 
}

/***************************************************************/
/***************************************************************/
/* CheckExistPursueType check if pursue type already exists */
func CheckExistPursueType(typePursueType string) (string, bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15 * time.Second)
	defer cancel()

	db := config.MongoConnection.Database("san_martin")
	collection := db.Collection("pursue_type")

	typePursueType = strings.ToUpper(typePursueType)
	condition := bson.M {
		"type": typePursueType,
	}

	var result models.PursueType

	err := collection.FindOne(ctx, condition).Decode(&result)
	if (result.Type != "") {
		return result.ID.Hex(), true, nil
	}

	return "", false, err
}

/***************************************************************/
/***************************************************************/
/* UpdatePursueTypeDB update the pursue type in the db */
func UpdatePursueTypeDB(p models.PursueType) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15 * time.Second)
	defer cancel()

	db := config.MongoConnection.Database("san_martin")
	collection := db.Collection("pursue_type")

	row := make(map[string]interface{})
	row["type"] = p.Type

	updateString := bson.M {
		"$set": row,
	}

	var idPursueType string
	idPursueType = p.ID.Hex()

	objID, _ := primitive.ObjectIDFromHex(idPursueType)

	filter := bson.M { "_id": bson.M { "$eq": objID }}

	_, err := collection.UpdateOne(ctx, filter, updateString)
	if err != nil {
		return false, err
	}

	return true, nil
}

/***************************************************************/
/***************************************************************/
/* DeletePursueTypeDB delete the pursue type from the db */
func DeletePursueTypeDB(IDPursueType string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 15 * time.Second)
	defer cancel()

	db := config.MongoConnection.Database("san_martin")
	collection := db.Collection("pursue_type")

	objID, _ := primitive.ObjectIDFromHex(IDPursueType)

	condition := bson.M {
		"_id": objID,
	}

	_, err := collection.DeleteOne(ctx, condition)
	return err
}

/***************************************************************/
/***************************************************************/
/* GetPursueTypeDB get the pursue type by id */
func GetPursueTypeDB(IDPursueType string) (models.PursueType, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15 * time.Second)
	defer cancel()

	db := config.MongoConnection.Database("san_martin")
	collection := db.Collection("pursue_type")

	objID, _ := primitive.ObjectIDFromHex(IDPursueType)

	condition := bson.M {
		"_id": objID,
	}

	var pursueType models.PursueType

	err := collection.FindOne(ctx, condition).Decode(&pursueType)
	return pursueType, err
}
