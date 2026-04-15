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

const teacherCollection = "teacher"

/* GetTeachersDB lista docentes */
func GetTeachersDB() ([]*models.Teacher, bool) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	collection := config.MongoConnection.Database("san_martin").Collection(teacherCollection)
	opts := options.Find().SetSort(bson.D{{Key: "name", Value: 1}})
	cur, err := collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		log.Println(err)
		return nil, false
	}
	defer cur.Close(ctx)
	var results []*models.Teacher
	for cur.Next(ctx) {
		var row models.Teacher
		if err := cur.Decode(&row); err != nil {
			return results, false
		}
		results = append(results, &row)
	}
	return results, true
}

/* InsertTeacherDB */
func InsertTeacherDB(t models.Teacher) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	collection := config.MongoConnection.Database("san_martin").Collection(teacherCollection)
	row := bson.M{
		"name":                t.Name,
		"email":               t.Email,
		"phone":               t.Phone,
		"dni":                 t.DNI,
		"address":             t.Address,
		"enseniaen":           t.EnseniaEn,
		"degreeids":           t.DegreeIDs,
		"titulohabilitanteid": t.TituloHabilitanteID,
		"modalidadid":         t.ModalidadID,
		"createdat":           t.CreatedAt,
		"updatedat":           t.UpdatedAt,
	}
	res, err := collection.InsertOne(ctx, row)
	if err != nil {
		return "Hubo un error al insertar el docente", err
	}
	oid, _ := res.InsertedID.(primitive.ObjectID)
	return oid.Hex(), nil
}

/* UpdateTeacherDB */
func UpdateTeacherDB(t models.Teacher) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	collection := config.MongoConnection.Database("san_martin").Collection(teacherCollection)
	filter := bson.M{"_id": t.ID}
	update := bson.M{"$set": bson.M{
		"name":                t.Name,
		"email":               t.Email,
		"phone":               t.Phone,
		"dni":                 t.DNI,
		"address":             t.Address,
		"enseniaen":           t.EnseniaEn,
		"degreeids":           t.DegreeIDs,
		"titulohabilitanteid": t.TituloHabilitanteID,
		"modalidadid":         t.ModalidadID,
		"updatedat":           t.UpdatedAt,
	}}
	_, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return false, err
	}
	return true, nil
}

/* FindTeacherByDNIDB igual que alumnos */
func FindTeacherByDNIDB(dni string, excludeID primitive.ObjectID) (models.Teacher, bool) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	collection := config.MongoConnection.Database("san_martin").Collection(teacherCollection)
	filter := bson.M{"dni": dni}
	if !excludeID.IsZero() {
		filter["_id"] = bson.M{"$ne": excludeID}
	}
	var te models.Teacher
	err := collection.FindOne(ctx, filter).Decode(&te)
	if err != nil {
		return te, false
	}
	return te, te.DNI != ""
}
