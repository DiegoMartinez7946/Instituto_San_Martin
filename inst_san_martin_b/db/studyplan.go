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
/* GetStudyPlansDB get the study plans from db */
func GetStudyPlansDB() ([]*models.StudyPlan, bool) {
	ctx, cancel := context.WithTimeout(context.Background(), 15 * time.Second)
	defer cancel()

	db := config.MongoConnection.Database("san_martin")
	collection := db.Collection("study_plan")

	var results []*models.StudyPlan

	condition := bson.M {  }
	optionsQuery := options.Find()
	optionsQuery.SetSort(bson.D {{ Key: "name", Value: 1}, { Key: "state", Value: 1}})

	studyPlans, err := collection.Find(ctx, condition, optionsQuery)
	if err != nil {
		log.Fatal(err.Error())
		return results, false
	}

	for studyPlans.Next(context.TODO()) {
		var row models.StudyPlan
		err := studyPlans.Decode(&row)
		if err != nil {
			return results, false
		}
		results = append(results, &row)
	}

	return results, true
}

/***************************************************************/
/***************************************************************/
/* InsertStudyPlanDB insert one study plan in db */
func InsertStudyPlanDB(s models.StudyPlan) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15 * time.Second)
	defer cancel()

	db := config.MongoConnection.Database("san_martin")
	collection := db.Collection("study_plan")

	row := bson.M {
		"name": s.Name,
		"code": s.Code,
		"state": s.State,
	}

	result, err := collection.InsertOne(ctx, row)
	if err != nil {
		return "Hubo un error al insertar el plan de estudio", err
	}
	
	objID, _ := result.InsertedID.(primitive.ObjectID)
	return objID.Hex(), nil 
}

/***************************************************************/
/***************************************************************/
/* CheckExistStudyPlan check if study plan already exists */
func CheckExistStudyPlan(code string) (string, bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15 * time.Second)
	defer cancel()

	db := config.MongoConnection.Database("san_martin")
	collection := db.Collection("study_plan")

	condition := bson.M {
		"code": code,
	}

	var result models.StudyPlan

	err := collection.FindOne(ctx, condition).Decode(&result)
	if (result.Name != "") {
		return result.ID.Hex(), true, nil
	}

	return "", false, err
}

/***************************************************************/
/***************************************************************/
/* UpdateStudyPlanDB update the study plan in the db */
func UpdateStudyPlanDB(s models.StudyPlan) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15 * time.Second)
	defer cancel()

	db := config.MongoConnection.Database("san_martin")
	collection := db.Collection("study_plan")

	row := make(map[string]interface{})
	row["name"] = s.Name

	updateString := bson.M {
		"$set": row,
	}

	var idStudyPlan string
	idStudyPlan = s.ID.Hex()

	objID, _ := primitive.ObjectIDFromHex(idStudyPlan)

	filter := bson.M { "_id": bson.M { "$eq": objID }}

	_, err := collection.UpdateOne(ctx, filter, updateString)
	if err != nil {
		return false, err
	}

	return true, nil
}

/***************************************************************/
/***************************************************************/
/* UpdateStateStudyPlanDB update the study plan in the db */
func UpdateStateStudyPlanDB(s models.StudyPlan) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15 * time.Second)
	defer cancel()

	db := config.MongoConnection.Database("san_martin")
	collection := db.Collection("study_plan")

	row := make(map[string]interface{})
	row["state"] = s.State

	updateString := bson.M {
		"$set": row,
	}

	var idStudyPlan string
	idStudyPlan = s.ID.Hex()

	objID, _ := primitive.ObjectIDFromHex(idStudyPlan)

	filter := bson.M { "_id": bson.M { "$eq": objID }}

	_, err := collection.UpdateOne(ctx, filter, updateString)
	if err != nil {
		return false, err
	}

	return true, nil
}
