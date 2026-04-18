package db

import (
	"context"
	"errors"
	"log"
	"strings"
	"time"

	"github.com/benjacifre10/san_martin_b/config"
	"github.com/benjacifre10/san_martin_b/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var ErrStudyPlanNotFound = errors.New("study plan not found")

const studyPlanCollection = "studyplan"

/* GetStudyPlansDB lista planes de estudio */
func GetStudyPlansDB() ([]*models.StudyPlan, bool) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	collection := config.WorkingDatabase().Collection(studyPlanCollection)
	var results []*models.StudyPlan
	opts := options.Find().SetSort(bson.D{{Key: "nombre", Value: 1}, {Key: "numeroresolucion", Value: 1}})
	cur, err := collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		log.Println("GetStudyPlansDB:", err.Error())
		return results, false
	}
	for cur.Next(ctx) {
		var row models.StudyPlan
		if err := cur.Decode(&row); err != nil {
			_ = cur.Close(ctx)
			return results, false
		}
		results = append(results, &row)
	}
	_ = cur.Close(ctx)
	return results, true
}

/* InsertStudyPlanDB inserta un plan */
func InsertStudyPlanDB(s models.StudyPlan) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	collection := config.WorkingDatabase().Collection(studyPlanCollection)
	row := bson.M{
		"nombre":           s.Nombre,
		"numeroresolucion": s.NumeroResolucion,
		"cohorte":          s.Cohorte,
		"extension":        s.Extension,
		"anioaprobacion":   s.AnioAprobacion,
		"aniocaducidad":    s.AnioCaducidad,
		"cargahoraria":     s.CargaHoraria,
		"active":           s.Active,
		"createdat":        s.CreatedAt,
		"updatedat":        s.UpdatedAt,
	}
	res, err := collection.InsertOne(ctx, row)
	if err != nil {
		return "Hubo un error al insertar el plan de estudio", err
	}
	oid, _ := res.InsertedID.(primitive.ObjectID)
	return oid.Hex(), nil
}

/* IsNumeroResolucionTakenByOther true si otro plan (no excludeID) usa ese número de resolución */
func IsNumeroResolucionTakenByOther(nr string, excludeID primitive.ObjectID) bool {
	nr = strings.TrimSpace(nr)
	if nr == "" {
		return false
	}
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	collection := config.WorkingDatabase().Collection(studyPlanCollection)
	var sp models.StudyPlan
	err := collection.FindOne(ctx, bson.M{"numeroresolucion": nr}).Decode(&sp)
	if err != nil || sp.ID.IsZero() {
		return false
	}
	if excludeID.IsZero() {
		return true
	}
	return sp.ID != excludeID
}

/* UpdateStudyPlanDB actualiza campos del plan */
func UpdateStudyPlanDB(s models.StudyPlan) (bool, error) {
	if s.ID.IsZero() {
		return false, ErrStudyPlanNotFound
	}
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	collection := config.WorkingDatabase().Collection(studyPlanCollection)
	row := bson.M{
		"nombre":           s.Nombre,
		"numeroresolucion": s.NumeroResolucion,
		"cohorte":          s.Cohorte,
		"extension":        s.Extension,
		"anioaprobacion":   s.AnioAprobacion,
		"aniocaducidad":    s.AnioCaducidad,
		"cargahoraria":     s.CargaHoraria,
		"updatedat":        s.UpdatedAt,
	}
	ur, err := collection.UpdateOne(ctx, bson.M{"_id": s.ID}, bson.M{"$set": row})
	if err != nil {
		return false, err
	}
	if ur.MatchedCount == 0 {
		return false, ErrStudyPlanNotFound
	}
	return true, nil
}

/* UpdateStudyPlanActiveDB solo activo/inactivo */
func UpdateStudyPlanActiveDB(id primitive.ObjectID, active bool) (bool, error) {
	if id.IsZero() {
		return false, nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	collection := config.WorkingDatabase().Collection(studyPlanCollection)
	ur, err := collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{
		"$set": bson.M{"active": active, "updatedat": time.Now()},
	})
	if err != nil {
		return false, err
	}
	if ur.MatchedCount == 0 {
		return false, ErrStudyPlanNotFound
	}
	return true, nil
}

/* StudyPlanDocumentExists existe documento con ese _id */
func StudyPlanDocumentExists(id primitive.ObjectID) bool {
	if id.IsZero() {
		return false
	}
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	collection := config.WorkingDatabase().Collection(studyPlanCollection)
	n, err := collection.CountDocuments(ctx, bson.M{"_id": id})
	return err == nil && n > 0
}

/* StudyPlanActiveByNumeroResolucion true si existe un plan activo con ese numeroresolucion */
func StudyPlanActiveByNumeroResolucion(nr string) bool {
	nr = strings.TrimSpace(nr)
	if nr == "" {
		return false
	}
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	collection := config.WorkingDatabase().Collection(studyPlanCollection)
	var sp models.StudyPlan
	err := collection.FindOne(ctx, bson.M{"numeroresolucion": nr, "active": true}).Decode(&sp)
	return err == nil && !sp.ID.IsZero()
}
