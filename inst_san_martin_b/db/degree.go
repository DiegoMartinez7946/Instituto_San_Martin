package db

import (
	"context"
	"log"
	"time"

	"github.com/benjacifre10/san_martin_b/config"
	"github.com/benjacifre10/san_martin_b/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

/***************************************************************/
/***************************************************************/
/* GetDegreesDB get the degrees from db */
func GetDegreesDB() ([]*models.Degree, bool) {
	ctx, cancel := context.WithTimeout(context.Background(), 15 * time.Second)
	defer cancel()

	db := config.MongoConnection.Database("san_martin")
	collection := db.Collection("degree")

	var results []*models.Degree

	condition := bson.M {  }
	optionsQuery := options.Find()
	optionsQuery.SetSort(bson.D {{ Key: "name", Value: 1}, { Key: "active", Value: 1}})

	degrees, err := collection.Find(ctx, condition, optionsQuery)
	if err != nil {
		log.Fatal(err.Error())
		return results, false
	}

	for degrees.Next(context.TODO()) {
		var row models.Degree
		err := degrees.Decode(&row)
		if err != nil {
			return results, false
		}
		results = append(results, &row)
	}

	return results, true
}

/***************************************************************/
/***************************************************************/
/* InsertDegreeDB insert one degree in db */
func InsertDegreeDB(d models.Degree) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15 * time.Second)
	defer cancel()

	db := config.MongoConnection.Database("san_martin")
	collection := db.Collection("degree")

	row := bson.M{
		"name":          d.Name,
		"nivel":         d.Nivel,
		"resolucionid":  d.ResolucionID,
		"active":        d.Active,
	}

	result, err := collection.InsertOne(ctx, row)
	if err != nil {
		return "Hubo un error al insertar la carrera", err
	}
	
	objID, _ := result.InsertedID.(primitive.ObjectID)
	return objID.Hex(), nil 
}

/***************************************************************/
/***************************************************************/
/* CheckExistDegree check if degree already exists */
func CheckExistDegree(nameDegree string) (string, bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15 * time.Second)
	defer cancel()

	db := config.MongoConnection.Database("san_martin")
	collection := db.Collection("degree")

	condition := bson.M {
		"name": nameDegree,
	}

	var result models.Degree

	err := collection.FindOne(ctx, condition).Decode(&result)
	if (result.Name != "") {
		return result.ID.Hex(), true, nil
	}

	return "", false, err
}

/* IsDegreeNameTakenByOther returns true if another degree (not excludeID) already uses this name */
func IsDegreeNameTakenByOther(nameDegree string, excludeID primitive.ObjectID) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	db := config.MongoConnection.Database("san_martin")
	collection := db.Collection("degree")

	var result models.Degree
	err := collection.FindOne(ctx, bson.M{"name": nameDegree}).Decode(&result)
	if err != nil || result.Name == "" || result.ID.IsZero() {
		return false
	}
	if excludeID.IsZero() {
		return true
	}
	return result.ID != excludeID
}

/* IsResolucionIDTakenByOther returns true if another degree already uses this resolucionid (trimmed exact match) */
func IsResolucionIDTakenByOther(resolucionID string, excludeID primitive.ObjectID) bool {
	if resolucionID == "" {
		return false
	}
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	db := config.MongoConnection.Database("san_martin")
	collection := db.Collection("degree")

	var result models.Degree
	err := collection.FindOne(ctx, bson.M{"resolucionid": resolucionID}).Decode(&result)
	if err != nil || result.ID.IsZero() {
		return false
	}
	if excludeID.IsZero() {
		return true
	}
	return result.ID != excludeID
}

/* GetDegreeByID returns one degree by id (any active flag) */
func GetDegreeByID(id primitive.ObjectID) (models.Degree, bool) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	db := config.MongoConnection.Database("san_martin")
	collection := db.Collection("degree")

	var d models.Degree
	err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&d)
	if err != nil || d.ID.IsZero() {
		return d, false
	}
	return d, true
}

/***************************************************************/
/***************************************************************/
/* UpdateDegreeDB update the degree in the db */
func UpdateDegreeDB(d models.Degree) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15 * time.Second)
	defer cancel()

	db := config.MongoConnection.Database("san_martin")
	collection := db.Collection("degree")

	row := make(map[string]interface{})
	row["name"] = d.Name
	row["nivel"] = d.Nivel
	row["resolucionid"] = d.ResolucionID

	updateString := bson.M {
		"$set": row,
	}

	var idDegree string
	idDegree = d.ID.Hex()

	objID, _ := primitive.ObjectIDFromHex(idDegree)

	filter := bson.M { "_id": bson.M { "$eq": objID }}

	_, err := collection.UpdateOne(ctx, filter, updateString)
	if err != nil {
		return false, err
	}

	return true, nil
}

/***************************************************************/
/***************************************************************/
/* UpdateStatusDegreeDB update the degree in the db */
func UpdateStatusDegreeDB(d models.Degree) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15 * time.Second)
	defer cancel()

	db := config.MongoConnection.Database("san_martin")
	collection := db.Collection("degree")

	row := make(map[string]interface{})
	row["active"] = d.Active

	updateString := bson.M {
		"$set": row,
	}

	var idDegree string
	idDegree = d.ID.Hex()

	objID, _ := primitive.ObjectIDFromHex(idDegree)

	filter := bson.M { "_id": bson.M { "$eq": objID }}

	_, err := collection.UpdateOne(ctx, filter, updateString)
	if err != nil {
		return false, err
	}

	return true, nil
}

/***************************************************************/
/***************************************************************/
/* ExistsDegreeByID returns true if a degree exists and is active */
func ExistsDegreeByID(id primitive.ObjectID) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 15 * time.Second)
	defer cancel()

	db := config.MongoConnection.Database("san_martin")
	collection := db.Collection("degree")

	var d models.Degree
	err := collection.FindOne(ctx, bson.M{"_id": id, "active": true}).Decode(&d)
	return err == nil && !d.ID.IsZero()
}

/* DegreeDocumentExists returns true if a degree document exists (any status) */
func DegreeDocumentExists(id primitive.ObjectID) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 15 * time.Second)
	defer cancel()

	db := config.MongoConnection.Database("san_martin")
	collection := db.Collection("degree")

	n, err := collection.CountDocuments(ctx, bson.M{"_id": id})
	return err == nil && n > 0
}
