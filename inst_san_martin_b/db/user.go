// Package db provides ...
package db

import (
	"context"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/benjacifre10/san_martin_b/config"
	"github.com/benjacifre10/san_martin_b/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const userCollection = "user"

var ensureUserActiveField sync.Once

func ensureUserActiveDefaults(ctx context.Context) {
	ensureUserActiveField.Do(func() {
		collection := config.MongoConnection.Database("san_martin").Collection(userCollection)
		_, err := collection.UpdateMany(ctx, bson.M{"active": bson.M{"$exists": false}}, bson.M{"$set": bson.M{"active": true}})
		if err != nil {
			log.Println("ensure user active:", err)
		}
	})
}

/***************************************************************/
/* InsertUserDB insert one user in db */
func InsertUserDB(u models.User) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	db := config.MongoConnection.Database("san_martin")
	collection := db.Collection(userCollection)

	result, err := collection.InsertOne(ctx, u)
	if err != nil {
		return "Hubo un error al insertar el usuario", err
	}

	ObjID, _ := result.InsertedID.(primitive.ObjectID)
	return ObjID.Hex(), nil
}

/***************************************************************/
/* CheckExistUser check if user already exists */
func CheckExistUser(email string) (models.User, bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	ensureUserActiveDefaults(ctx)

	db := config.MongoConnection.Database("san_martin")
	collection := db.Collection(userCollection)

	condition := bson.M{
		"email": email,
	}

	var result models.User

	err := collection.FindOne(ctx, condition).Decode(&result)
	if result.Email != "" {
		return result, true, nil
	}

	return result, false, err
}

/***************************************************************/
/* GetUserDB get one user from db */
func GetUserDB(email string) (models.UserResponse, bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	ensureUserActiveDefaults(ctx)

	db := config.MongoConnection.Database("san_martin")
	collection := db.Collection(userCollection)

	condition := make([]bson.M, 0)
	condition = append(condition, bson.M{
		"$project": bson.M{
			"_id":      1,
			"usertype": bson.M{"$toObjectId": "$usertype"},
			"email":    "$email",
		},
	})
	condition = append(condition, bson.M{
		"$lookup": bson.M{
			"from":         "role",
			"localField":   "usertype",
			"foreignField": "_id",
			"as":           "role",
		}})
	condition = append(condition, bson.M{"$unwind": "$role"})
	condition = append(condition, bson.M{"$match": bson.M{"email": email}})
	condition = append(condition, bson.M{
		"$project": bson.M{
			"_id":   1,
			"email": "$email",
			"role":  "$role.type",
		},
	})

	cur, err := collection.Aggregate(ctx, condition)
	var result []models.UserResponse

	err = cur.All(ctx, &result)
	if err != nil || len(result) == 0 {
		return models.UserResponse{}, false, err
	}
	return result[0], true, nil
}

/***************************************************************/
/* GetUsersDB get the users from db */
func GetUsersDB(roleType string) ([]models.UserListRow, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	ensureUserActiveDefaults(ctx)

	db := config.MongoConnection.Database("san_martin")
	collection := db.Collection(userCollection)

	condition := make([]bson.M, 0)

	condition = append(condition, bson.M{
		"$project": bson.M{
			"_id":      1,
			"email":    1,
			"usertype": bson.M{"$toObjectId": "$usertype"},
			"name":     bson.M{"$ifNull": bson.A{"$name", ""}},
			"dni":      bson.M{"$ifNull": bson.A{"$dni", ""}},
			"address":  bson.M{"$ifNull": bson.A{"$address", ""}},
			"phone":    bson.M{"$ifNull": bson.A{"$phone", ""}},
			"active":   bson.M{"$ifNull": bson.A{"$active", true}},
			"shiftidOne": bson.M{"$ifNull": bson.A{"$shiftid", ""}},
			"shiftidsArr": bson.M{"$ifNull": bson.A{"$shiftids", bson.A{}}},
		},
	})
	condition = append(condition, bson.M{
		"$lookup": bson.M{
			"from":         "role",
			"localField":   "usertype",
			"foreignField": "_id",
			"as":           "role",
		}})
	condition = append(condition, bson.M{"$match": bson.M{"role.type": roleType}})
	condition = append(condition, bson.M{"$unwind": "$role"})

	condition = append(condition, bson.M{"$addFields": bson.M{
		"shiftIdsNorm": bson.M{
			"$cond": bson.A{
				bson.M{"$gt": bson.A{bson.M{"$size": "$shiftidsArr"}, 0}},
				"$shiftidsArr",
				bson.M{
					"$cond": bson.A{
						bson.M{"$ne": bson.A{"$shiftidOne", ""}},
						bson.A{"$shiftidOne"},
						bson.A{},
					},
				},
			},
		},
	}})

	condition = append(condition, bson.M{
		"$lookup": bson.M{
			"from": "shift",
			"let":  bson.M{"ids": "$shiftIdsNorm"},
			"pipeline": []bson.M{
				{"$match": bson.M{
					"$expr": bson.M{
						"$in": bson.A{bson.M{"$toString": "$_id"}, "$$ids"},
					},
				}},
			},
			"as": "shiftsdoc",
		},
	})

	condition = append(condition, bson.M{"$addFields": bson.M{
		"shiftTypeJoined": bson.M{
			"$reduce": bson.M{
				"input":        "$shiftsdoc",
				"initialValue": "",
				"in": bson.M{
					"$cond": bson.A{
						bson.M{"$eq": bson.A{"$$value", ""}},
						"$$this.type",
						bson.M{"$concat": bson.A{"$$value", ", ", "$$this.type"}},
					},
				},
			},
		},
	}})

	condition = append(condition, bson.M{
		"$project": bson.M{
			"id":        bson.M{"$toString": "$_id"},
			"email":     1,
			"role":      "$role.type",
			"name":      1,
			"dni":       1,
			"address":   1,
			"phone":     1,
			"active":    1,
			"shiftIds":  "$shiftIdsNorm",
			"shiftType": bson.M{"$ifNull": bson.A{"$shiftTypeJoined", ""}},
		},
	})

	cur, err := collection.Aggregate(ctx, condition)
	var result []models.UserListRow

	err = cur.All(ctx, &result)
	if err != nil {
		return result, 400, err
	}
	return result, 200, nil
}

/***************************************************************/
/* FindOtherUserByDNIDB true si otro usuario ya usa ese DNI */
func FindOtherUserByDNIDB(dni string, excludeID primitive.ObjectID) (bool, error) {
	d := strings.TrimSpace(dni)
	if d == "" {
		return false, nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	collection := config.MongoConnection.Database("san_martin").Collection(userCollection)
	filter := bson.M{"dni": d}
	if !excludeID.IsZero() {
		filter["_id"] = bson.M{"$ne": excludeID}
	}
	n, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return false, err
	}
	return n > 0, nil
}

/***************************************************************/
/* FindUserByIDDB obtiene usuario por id */
func FindUserByIDDB(id primitive.ObjectID) (models.User, bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	ensureUserActiveDefaults(ctx)

	var u models.User
	err := config.MongoConnection.Database("san_martin").Collection(userCollection).FindOne(ctx, bson.M{"_id": id}).Decode(&u)
	if err == mongo.ErrNoDocuments {
		return u, false, nil
	}
	if err != nil {
		return u, false, err
	}
	if len(u.ShiftIDs) == 0 && strings.TrimSpace(u.ShiftID) != "" {
		u.ShiftIDs = []string{strings.TrimSpace(u.ShiftID)}
	}
	return u, true, nil
}

/***************************************************************/
/* UpdateUserProfileDB actualiza datos de perfil (sin password) */
func UpdateUserProfileDB(id primitive.ObjectID, name, dni, address, phone, email string, shiftIDs []string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	collection := config.MongoConnection.Database("san_martin").Collection(userCollection)
	setDoc := bson.M{
		"name":      name,
		"dni":       dni,
		"address":   address,
		"phone":     phone,
		"email":     email,
		"shiftids":  shiftIDs,
		"updatedat": time.Now(),
	}
	_, err := collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{
		"$set":   setDoc,
		"$unset": bson.M{"shiftid": ""},
	})
	if err != nil {
		return false, err
	}
	return true, nil
}

/***************************************************************/
/* UpdateUserPasswordByIDDB actualiza password por id */
func UpdateUserPasswordByIDDB(id primitive.ObjectID, hashedPassword string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	collection := config.MongoConnection.Database("san_martin").Collection(userCollection)
	_, err := collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{
		"$set": bson.M{
			"password":  hashedPassword,
			"updatedat": time.Now(),
		},
	})
	if err != nil {
		return false, err
	}
	return true, nil
}

/***************************************************************/
/* UpdateUserActiveDB activo / inactivo */
func UpdateUserActiveDB(id primitive.ObjectID, active bool) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	collection := config.MongoConnection.Database("san_martin").Collection(userCollection)
	_, err := collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{
		"$set": bson.M{
			"active":    active,
			"updatedat": time.Now(),
		},
	})
	if err != nil {
		return false, err
	}
	return true, nil
}

/***************************************************************/
/* CheckExistUserByEmailOtherThanID email tomado por otro usuario */
func CheckExistUserByEmailOtherThanID(email string, excludeID primitive.ObjectID) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	filter := bson.M{"email": strings.TrimSpace(email)}
	if !excludeID.IsZero() {
		filter["_id"] = bson.M{"$ne": excludeID}
	}
	n, err := config.MongoConnection.Database("san_martin").Collection(userCollection).CountDocuments(ctx, filter)
	if err != nil {
		return false, err
	}
	return n > 0, nil
}

/***************************************************************/
/* UpdateUserRoleDB actualiza el rol (usertype = ObjectId del rol en hex) */
func UpdateUserRoleDB(userID primitive.ObjectID, roleIDHex string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	collection := config.MongoConnection.Database("san_martin").Collection(userCollection)
	_, err := collection.UpdateOne(ctx, bson.M{"_id": userID}, bson.M{
		"$set": bson.M{
			"usertype":  strings.TrimSpace(roleIDHex),
			"updatedat": time.Now(),
		},
	})
	if err != nil {
		return false, err
	}
	return true, nil
}

/***************************************************************/
/* ChangePasswordDB update the password in the db */
func ChangePasswordDB(u models.User) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	db := config.MongoConnection.Database("san_martin")
	collection := db.Collection(userCollection)

	row := make(map[string]interface{})
	row["password"] = u.Password
	row["updatedat"] = time.Now()

	updateString := bson.M{
		"$set": row,
	}

	filter := bson.M{"email": bson.M{"$eq": strings.TrimSpace(u.Email)}}

	res, err := collection.UpdateOne(ctx, filter, updateString)
	if err != nil {
		return false, err
	}
	if res.MatchedCount == 0 {
		return false, nil
	}

	return true, nil
}
