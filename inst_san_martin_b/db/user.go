// Package db provides ...
package db

import (
	"context"
	"time"

	"github.com/benjacifre10/san_martin_b/config"
	"github.com/benjacifre10/san_martin_b/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

/* InsertUserDB insert one user in db */
func InsertUserDB(u models.User) (string, bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15 * time.Second)
	defer cancel()

	db := config.MongoConnection.Database("san_martin")
	collection := db.Collection("user")

	result, err := collection.InsertOne(ctx, u)
	if err != nil {
		return "", false, err
	}

	ObjID, _ := result.InsertedID.(primitive.ObjectID)
	return ObjID.Hex(), true, nil
}

/* CheckExistUser check if user already exists */
func CheckExistUser(email string) (models.User, bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15 * time.Second)
	defer cancel()

	db := config.MongoConnection.Database("san_martin")
	collection := db.Collection("user")

	condition := bson.M {
		"email": email,
	}

	var result models.User

	err := collection.FindOne(ctx, condition).Decode(&result)
	if (result.Email != "") {
		return result, true, nil
	}

	return result, false, err
}

/* GetUserDB get one user from db */
func GetUserDB(email string) (models.UserResponse, bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15 * time.Second)
	defer cancel()

	db := config.MongoConnection.Database("san_martin")
	collection := db.Collection("user")

	condition := make([]bson.M, 0)
	condition = append(condition, bson.M {
		"$project": bson.M { 
			"usertype": bson.M { "$toObjectId": "$usertype" },
			"email": "$email",
		},
	})
	condition = append(condition, bson.M {
		"$lookup": bson.M {
			"from": "role",
			"localField": "usertype",
			"foreignField": "_id",
			"as": "role",
	}})
	condition = append(condition, bson.M { "$unwind": "$role" })
	condition = append(condition, bson.M { "$match": bson.M { "email": email }})

	cur, err := collection.Aggregate(ctx, condition)
	var result []models.UserResponse

	err = cur.All(ctx, &result)
	if err != nil {
		return result[0], false, err
  }
	return result[0], true, nil
}

/* GetUsersDB get the users from db */
func GetUsersDB(roleType string) ([]models.UserResponse, bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15 * time.Second)
	defer cancel()

	db := config.MongoConnection.Database("san_martin")
	collection := db.Collection("user")

	condition := make([]bson.M, 0)
	condition = append(condition, bson.M {
		"$project": bson.M { 
			"usertype": bson.M { "$toObjectId": "$usertype" },
			"email": "$email",
		},
	})
	condition = append(condition, bson.M {
		"$lookup": bson.M {
			"from": "role",
			"localField": "usertype",
			"foreignField": "_id",
			"as": "role",
	}})
	condition = append(condition, bson.M { "$unwind": "$role" })
	condition = append(condition, bson.M { "$match": bson.M { "role.type": roleType } })

	cur, err := collection.Aggregate(ctx, condition)
	var result []models.UserResponse

	err = cur.All(ctx, &result)
	if err != nil {
		return result, false, err
  }
	return result, true, nil
}

/* ChangePasswordDB update the password in the db */
func ChangePasswordDB(u models.User) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15 * time.Second)
	defer cancel()

	db := config.MongoConnection.Database("san_martin")
	collection := db.Collection("user")

	row := make(map[string]interface{})
	row["password"] = u.Password

	updateString := bson.M {
		"$set": row,
	}

//	var idRole string
//	idRole = r.ID.Hex()

//	objID, _ := primitive.ObjectIDFromHex(idRole)

	filter := bson.M { "email": bson.M { "$eq": u.Email }}

	_, err := collection.UpdateOne(ctx, filter, updateString)
	if err != nil {
		return false, err
	}

	return true, nil
}
