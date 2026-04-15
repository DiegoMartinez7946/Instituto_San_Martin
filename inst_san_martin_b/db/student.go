package db

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/benjacifre10/san_martin_b/config"
	"github.com/benjacifre10/san_martin_b/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var ensureStudentActiveField sync.Once

func mergeStudentDegreeIDsFromLegacy(s *models.Student) {
	if s == nil {
		return
	}
	if len(s.DegreeIDs) == 0 && len(s.DegreeIDsLegacy) > 0 {
		s.DegreeIDs = append([]primitive.ObjectID(nil), s.DegreeIDsLegacy...)
	}
}

func ensureStudentActiveDefaults(ctx context.Context) {
	ensureStudentActiveField.Do(func() {
		db := config.MongoConnection.Database("san_martin")
		collection := db.Collection("student")
		_, err := collection.UpdateMany(ctx, bson.M{"active": bson.M{"$exists": false}}, bson.M{"$set": bson.M{"active": true}})
		if err != nil {
			log.Println("ensure student active:", err)
		}
		// Documentos viejos sin claves: añade campos vacíos sin pisar valores existentes.
		if _, e2 := collection.UpdateMany(ctx, bson.M{"modalidad": bson.M{"$exists": false}}, bson.M{"$set": bson.M{"modalidad": ""}}); e2 != nil {
			log.Println("ensure student modalidad:", e2)
		}
		if _, e3 := collection.UpdateMany(ctx, bson.M{"condicion": bson.M{"$exists": false}}, bson.M{"$set": bson.M{"condicion": ""}}); e3 != nil {
			log.Println("ensure student condicion:", e3)
		}
	})
}

/* GetStudentsDB get all students */
func GetStudentsDB() ([]*models.Student, bool) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	ensureStudentActiveDefaults(ctx)

	db := config.MongoConnection.Database("san_martin")
	collection := db.Collection("student")

	var results []*models.Student

	opts := options.Find().SetSort(bson.D{{Key: "name", Value: 1}})
	cur, err := collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		log.Println(err.Error())
		return results, false
	}

	for cur.Next(ctx) {
		var row models.Student
		if err := cur.Decode(&row); err != nil {
			return results, false
		}
		mergeStudentDegreeIDsFromLegacy(&row)
		results = append(results, &row)
	}
	_ = cur.Close(ctx)
	return results, true
}

/* InsertStudentDB insert one student */
func InsertStudentDB(s models.Student) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	db := config.MongoConnection.Database("san_martin")
	collection := db.Collection("student")

	row := bson.M{
		"name":          s.Name,
		"email":         s.Email,
		"phone":         s.Phone,
		"dni":           s.DNI,
		"address":       s.Address,
		"nivelaprobado": s.NivelAprobado,
		"modalidad":     s.Modalidad,
		"condicion":     s.Condicion,
		"degreeids":     s.DegreeIDs,
		"active":        s.Active,
		"createdat":     s.CreatedAt,
		"updatedat":     s.UpdatedAt,
	}

	result, err := collection.InsertOne(ctx, row)
	if err != nil {
		return "Hubo un error al insertar el alumno", err
	}

	objID, _ := result.InsertedID.(primitive.ObjectID)
	return objID.Hex(), nil
}

/* UpdateStudentDB update student */
func UpdateStudentDB(s models.Student) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	db := config.MongoConnection.Database("san_martin")
	collection := db.Collection("student")

	filter := bson.M{"_id": s.ID}
	update := bson.M{
		"$set": bson.M{
			"name":          s.Name,
			"email":         s.Email,
			"phone":         s.Phone,
			"dni":           s.DNI,
			"address":       s.Address,
			"nivelaprobado": s.NivelAprobado,
			"modalidad":     s.Modalidad,
			"condicion":     s.Condicion,
			"degreeids":     s.DegreeIDs,
			"active":        s.Active,
			"updatedat":     s.UpdatedAt,
		},
		// Clave legacy en camelCase; el driver Go usa "degreeids".
		"$unset": bson.M{"degreeIds": ""},
	}

	_, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return false, err
	}
	return true, nil
}

/* FindStudentByEmailInsensitiveDB busca alumno por email (coincidencia sin distinguir mayúsculas) */
func FindStudentByEmailInsensitiveDB(email string) (models.Student, bool) {
	em := strings.TrimSpace(email)
	if em == "" {
		return models.Student{}, false
	}
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	db := config.MongoConnection.Database("san_martin")
	collection := db.Collection("student")

	pattern := "^" + regexp.QuoteMeta(em) + "$"
	filter := bson.M{"email": bson.M{"$regex": pattern, "$options": "i"}}

	var st models.Student
	err := collection.FindOne(ctx, filter).Decode(&st)
	if err != nil || st.ID.IsZero() {
		return models.Student{}, false
	}
	mergeStudentDegreeIDsFromLegacy(&st)
	return st, true
}

/* FindStudentByIDDB carga una ficha por _id */
func FindStudentByIDDB(id primitive.ObjectID) (models.Student, bool) {
	if id.IsZero() {
		return models.Student{}, false
	}
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	db := config.MongoConnection.Database("san_martin")
	collection := db.Collection("student")

	var st models.Student
	err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&st)
	if err != nil {
		return models.Student{}, false
	}
	mergeStudentDegreeIDsFromLegacy(&st)
	if st.ID.IsZero() {
		return models.Student{}, false
	}
	return st, true
}

/* FindStudentForUserSync ficha alumno por email (insensible a mayúsculas) o por DNI */
func FindStudentForUserSync(email, dni string) (models.Student, bool) {
	if st, ok := FindStudentByEmailInsensitiveDB(email); ok {
		return st, true
	}
	d := strings.TrimSpace(dni)
	if d == "" {
		return models.Student{}, false
	}
	return FindStudentByDNIDB(d, primitive.NilObjectID)
}

/* UpdateStudentEnrollmentFieldsDB actualiza carrera y modalidad/condición (sincronización desde usuario) */
func UpdateStudentEnrollmentFieldsDB(id primitive.ObjectID, degreeIDs []primitive.ObjectID, modalidad, condicion string) (bool, error) {
	if id.IsZero() {
		return false, nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	db := config.MongoConnection.Database("san_martin")
	collection := db.Collection("student")

	if degreeIDs == nil {
		degreeIDs = []primitive.ObjectID{}
	}

	res, err := collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{
		"$set": bson.M{
			"degreeids": degreeIDs,
			"modalidad": modalidad,
			"condicion": condicion,
			"updatedat": time.Now(),
		},
		"$unset": bson.M{"degreeIds": ""},
	})
	if err != nil {
		return false, err
	}
	if res.MatchedCount == 0 {
		return false, fmt.Errorf("no se encontro documento student con ese id")
	}
	return true, nil
}

/* FindStudentByDNIDB returns student if DNI exists (optional exclude id for update) */
func FindStudentByDNIDB(dni string, excludeID primitive.ObjectID) (models.Student, bool) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	db := config.MongoConnection.Database("san_martin")
	collection := db.Collection("student")

	filter := bson.M{"dni": dni}
	if !excludeID.IsZero() {
		filter["_id"] = bson.M{"$ne": excludeID}
	}

	var st models.Student
	err := collection.FindOne(ctx, filter).Decode(&st)
	if err != nil {
		return st, false
	}
	mergeStudentDegreeIDsFromLegacy(&st)
	return st, st.DNI != ""
}

/* UpdateStudentActiveDB solo estado activo/inactivo */
func UpdateStudentActiveDB(id primitive.ObjectID, active bool) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	db := config.MongoConnection.Database("san_martin")
	collection := db.Collection("student")

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

/* ClearStudentModeFieldsDB deja modalidad/condicion en blanco */
func ClearStudentModeFieldsDB(id primitive.ObjectID) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	db := config.MongoConnection.Database("san_martin")
	collection := db.Collection("student")

	_, err := collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{
		"$set": bson.M{
			"modalidad": "",
			"condicion": "",
			"updatedat": time.Now(),
		},
	})
	if err != nil {
		return false, err
	}
	return true, nil
}
