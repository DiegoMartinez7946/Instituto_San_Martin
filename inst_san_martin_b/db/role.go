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

/* GetRolesDB get the roles from db */
func GetRolesDB() ([]*models.Role, bool) {
	ctx, cancel := context.WithTimeout(context.Background(), 15 * time.Second)
	defer cancel()

	db := config.MongoConnection.Database("san_martin")
	collection := db.Collection("role")

	var results []*models.Role

	condition := bson.M {  }
	optionsQuery := options.Find()
	optionsQuery.SetSort(bson.D {{ Key: "type", Value: -1}})

	roles, err := collection.Find(ctx, condition, optionsQuery)
	if err != nil {
		log.Fatal(err.Error())
		return results, false
	}

	for roles.Next(context.TODO()) {
		var row models.Role
		err := roles.Decode(&row)
		if err != nil {
			return results, false
		}
		results = append(results, &row)
	}

	return results, true
}

/* InsertRoleDB insert one role in db */
func InsertRoleDB(r models.Role) (string, bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15 * time.Second)
	defer cancel()

	db := config.MongoConnection.Database("san_martin")
	collection := db.Collection("role")

	row := bson.M {
		"type": r.Type,
	}

	result, err := collection.InsertOne(ctx, row)
	if err != nil {
		return "", false, err
	}
	
	objID, _ := result.InsertedID.(primitive.ObjectID)
	return objID.Hex(), true, nil 
}

/* CheckExistRole check if role already exists */
func CheckExistRole(typeRol string) (string, bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15 * time.Second)
	defer cancel()

	db := config.MongoConnection.Database("san_martin")
	collection := db.Collection("role")

	typeRol = strings.ToUpper(typeRol)
	condition := bson.M {
		"type": typeRol,
	}

	var result models.Role

	err := collection.FindOne(ctx, condition).Decode(&result)
	if (result.Type != "") {
		return result.ID.Hex(), true, nil
	}

	return "", false, err
}

/* UpdateRoleDB update the role in the db */
func UpdateRoleDB(r models.Role) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15 * time.Second)
	defer cancel()

	db := config.MongoConnection.Database("san_martin")
	collection := db.Collection("role")

	row := make(map[string]interface{})
	row["type"] = r.Type

	updateString := bson.M {
		"$set": row,
	}

	var idRole string
	idRole = r.ID.Hex()

	objID, _ := primitive.ObjectIDFromHex(idRole)

	filter := bson.M { "_id": bson.M { "$eq": objID }}

	_, err := collection.UpdateOne(ctx, filter, updateString)
	if err != nil {
		return false, err
	}

	return true, nil
}

/* DeleteRoleDB delete the user role from the db */
func DeleteRoleDB(IDRole string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 15 * time.Second)
	defer cancel()

	db := config.MongoConnection.Database("san_martin")
	collection := db.Collection("role")

	objID, _ := primitive.ObjectIDFromHex(IDRole)

	condition := bson.M {
		"_id": objID,
	}

	_, err := collection.DeleteOne(ctx, condition)
	return err
}
