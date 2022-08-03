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

/***************************************************************/
/***************************************************************/
/* InsertUserDB insert one user in db */
func InsertUserDB(u models.User) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15 * time.Second)
	defer cancel()

	db := config.MongoConnection.Database("san_martin")
	collection := db.Collection("user")

	result, err := collection.InsertOne(ctx, u)
	if err != nil {
		return "Hubo un error al insertar el usuario", err
	}

	ObjID, _ := result.InsertedID.(primitive.ObjectID)
	return ObjID.Hex(), nil
}

/***************************************************************/
/***************************************************************/
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

/***************************************************************/
/***************************************************************/
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
	condition = append(condition, bson.M {
		"$project": bson.M { 
			"email": "$email",
			"role": "$role.type",
		},
	})

	cur, err := collection.Aggregate(ctx, condition)
	var result []models.UserResponse

	err = cur.All(ctx, &result)
	if err != nil {
		return result[0], false, err
  }
	return result[0], true, nil
}

/***************************************************************/
/***************************************************************/
/* GetUsersDB get the users from db */
func GetUsersDB(roleType string) ([]models.UserResponse, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15 * time.Second)
	defer cancel()

	db := config.MongoConnection.Database("san_martin")
	collection := db.Collection("user")

	condition := make([]bson.M, 0)

	// project me sirve tanto para dejar afuera a algunos campos
	// como tambien para calcular y mostrar otros
	condition = append(condition, bson.M {
		"$project": bson.M { 
			"usertype": bson.M { "$toObjectId": "$usertype" },
			"email": "$email",
		},
	})
	// aca conecto las dos tablas
	condition = append(condition, bson.M {
		"$lookup": bson.M {
			"from": "role",
			"localField": "usertype",
			"foreignField": "_id",
			"as": "role",
	}})
	// aca agrego la condicion para que el role type sea igual al del parametro de entrada
	condition = append(condition, bson.M { "$match": bson.M { "role.type": roleType } })
	// con esto no aplico como si fuera una coleccion de users con arrays de role
	// sino simplemente una coleccion de users con un rol por cada user
	condition = append(condition, bson.M { "$unwind": "$role" })

	condition = append(condition, bson.M {
		"$project": bson.M { 
			"email": "$email",
			"role": "$role.type",
		},
	})
	
	cur, err := collection.Aggregate(ctx, condition)
	var result []models.UserResponse

	err = cur.All(ctx, &result)
	if err != nil {
		return result, 400, err
  }
	return result, 200, nil
}

/***************************************************************/
/***************************************************************/
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

	filter := bson.M { "email": bson.M { "$eq": u.Email }}

	_, err := collection.UpdateOne(ctx, filter, updateString)
	if err != nil {
		return false, err
	}

	return true, nil
}
