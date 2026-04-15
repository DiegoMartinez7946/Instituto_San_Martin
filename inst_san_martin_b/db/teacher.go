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
	"go.mongodb.org/mongo-driver/mongo/options"
)

const teacherCollection = "teacher"

var ensureTeacherActiveField sync.Once

func ensureTeacherActiveDefaults(ctx context.Context) {
	ensureTeacherActiveField.Do(func() {
		collection := config.MongoConnection.Database("san_martin").Collection(teacherCollection)
		_, err := collection.UpdateMany(ctx, bson.M{"active": bson.M{"$exists": false}}, bson.M{"$set": bson.M{"active": true}})
		if err != nil {
			log.Println("ensure teacher active:", err)
		}
	})
}

/* GetTeachersDB lista docentes */
func GetTeachersDB() ([]*models.Teacher, bool) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	ensureTeacherActiveDefaults(ctx)
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
	careers := bson.A{}
	for _, c := range t.Careers {
		careers = append(careers, bson.M{
			"degreeid":            c.DegreeID,
			"titulohabilitanteid": c.TituloHabilitanteID,
			"modalidadid":         c.ModalidadID,
		})
	}
	row := bson.M{
		"name":      t.Name,
		"email":     t.Email,
		"phone":     t.Phone,
		"dni":       t.DNI,
		"address":   t.Address,
		"enseniaen": t.EnseniaEn,
		"careers":   careers,
		"active":    t.Active,
		"createdat": t.CreatedAt,
		"updatedat": t.UpdatedAt,
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
	careers := bson.A{}
	for _, c := range t.Careers {
		careers = append(careers, bson.M{
			"degreeid":            c.DegreeID,
			"titulohabilitanteid": c.TituloHabilitanteID,
			"modalidadid":         c.ModalidadID,
		})
	}
	filter := bson.M{"_id": t.ID}
	_, err := collection.UpdateOne(ctx, filter, bson.M{
		"$set": bson.M{
			"name":      t.Name,
			"email":     t.Email,
			"phone":     t.Phone,
			"dni":       t.DNI,
			"address":   t.Address,
			"enseniaen": t.EnseniaEn,
			"careers":   careers,
			"active":    t.Active,
			"updatedat": t.UpdatedAt,
		},
		"$unset": bson.M{
			"degreeids":           "",
			"titulohabilitanteid": "",
			"modalidadid":         "",
		},
	})
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

/* GetTeacherByID un docente por id */
func GetTeacherByID(id primitive.ObjectID) (models.Teacher, bool) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	collection := config.MongoConnection.Database("san_martin").Collection(teacherCollection)
	var t models.Teacher
	err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&t)
	if err != nil || t.ID.IsZero() {
		return t, false
	}
	return t, true
}

/* UpdateTeacherActiveDB solo estado activo/inactivo */
func UpdateTeacherActiveDB(id primitive.ObjectID, active bool) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	collection := config.MongoConnection.Database("san_martin").Collection(teacherCollection)
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

/* tituloIDModalidadIDByCodigo busca ids en catálogos (deben existir tras seed) */
func tituloIDByCodigo(codigo string) (primitive.ObjectID, bool) {
	list, ok := GetTitulosHabilitantesDB()
	if !ok {
		return primitive.ObjectID{}, false
	}
	u := strings.ToUpper(strings.TrimSpace(codigo))
	for _, x := range list {
		if strings.ToUpper(strings.TrimSpace(x.Codigo)) == u {
			return x.ID, true
		}
	}
	return primitive.ObjectID{}, false
}

func modalidadIDByCodigo(codigo string) (primitive.ObjectID, bool) {
	list, ok := GetModalidadesDB()
	if !ok {
		return primitive.ObjectID{}, false
	}
	u := strings.ToUpper(strings.TrimSpace(codigo))
	for _, x := range list {
		if strings.ToUpper(strings.TrimSpace(x.Codigo)) == u {
			return x.ID, true
		}
	}
	return primitive.ObjectID{}, false
}

/* FillTeacherCareersFromLegacy si el documento solo tenía degreeids */
func FillTeacherCareersFromLegacy(t *models.Teacher) {
	if len(t.Careers) > 0 || len(t.LegacyDegreeIDs) == 0 {
		return
	}
	noID, ok1 := tituloIDByCodigo("NO")
	provID, ok2 := modalidadIDByCodigo("PROVISIONAL")
	if !ok1 || !ok2 {
		return
	}
	for _, did := range t.LegacyDegreeIDs {
		if did.IsZero() {
			continue
		}
		t.Careers = append(t.Careers, models.TeacherCareerAssignment{
			DegreeID:            did,
			TituloHabilitanteID: noID,
			ModalidadID:         provID,
		})
	}
}
