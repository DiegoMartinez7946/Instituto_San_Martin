package config

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/benjacifre10/san_martin_b/models"
	"github.com/benjacifre10/san_martin_b/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

/* MongoConnection is a export variable to access anywhere */
var MongoConnection = ConnectDB()

var dbName = utils.GoDotEnvValue("DB_NAME")

/* WorkingDatabase base de datos de trabajo (DB_NAME en .env; por defecto san_martin). */
func WorkingDatabase() *mongo.Database {
	n := strings.TrimSpace(dbName)
	if n == "" {
		n = "san_martin"
	}
	return MongoConnection.Database(n)
}
var dbUser = utils.GoDotEnvValue("DB_USER")
var dbPass = utils.GoDotEnvValue("DB_PASSWORD")
var dbCluster = utils.GoDotEnvValue("DB_CLUSTER")

func getConnectionString() string {

	var connectionString strings.Builder  

	connectionString.WriteString("mongodb+srv://")
	connectionString.WriteString(dbUser + ":" + dbPass)
	connectionString.WriteString("@" + dbCluster + "/" + dbName)
	connectionString.WriteString("?retryWrites=true&w=majority")

	return connectionString.String()
}

var clientOptions = options.Client().ApplyURI(getConnectionString())

/* migratePursueTypeToAttendanceMode renombra la colección antigua pursue_type a attendance_mode (una sola vez). */
func migratePursueTypeToAttendanceMode(client *mongo.Client) {
	name := strings.TrimSpace(dbName)
	if name == "" {
		name = "san_martin"
	}
	ctx, cancel := context.WithTimeout(context.Background(), 25*time.Second)
	defer cancel()

	db := client.Database(name)
	names, err := db.ListCollectionNames(ctx, bson.M{})
	if err != nil {
		log.Printf("migración pursue_type: listar colecciones: %v", err)
		return
	}
	var hasOld, hasNew bool
	for _, n := range names {
		switch n {
		case "pursue_type":
			hasOld = true
		case "attendance_mode":
			hasNew = true
		}
	}
	if !hasOld {
		return
	}
	if hasNew {
		log.Printf("migración pursue_type: existen pursue_type y attendance_mode; elimine o fusione pursue_type manualmente en MongoDB")
		return
	}
	from := name + ".pursue_type"
	to := name + ".attendance_mode"
	err = client.Database("admin").RunCommand(ctx, bson.D{
		{Key: "renameCollection", Value: from},
		{Key: "to", Value: to},
		{Key: "dropTarget", Value: false},
	}).Err()
	if err != nil {
		log.Printf("migración pursue_type: renameCollection falló (%v). En Atlas: renombre la colección pursue_type a attendance_mode", err)
		return
	}
	log.Printf("migración pursue_type: colección renombrada a %s", to)
}

/* migrateJerarquiaToLevel renombra jerarquia → level */
func migrateJerarquiaToLevel(client *mongo.Client) {
	name := strings.TrimSpace(dbName)
	if name == "" {
		name = "san_martin"
	}
	ctx, cancel := context.WithTimeout(context.Background(), 25*time.Second)
	defer cancel()

	db := client.Database(name)
	names, err := db.ListCollectionNames(ctx, bson.M{})
	if err != nil {
		log.Printf("migración jerarquia: listar colecciones: %v", err)
		return
	}
	var hasOld, hasNew bool
	for _, n := range names {
		switch n {
		case "jerarquia":
			hasOld = true
		case "level":
			hasNew = true
		}
	}
	if !hasOld {
		return
	}
	if hasNew {
		log.Printf("migración jerarquia: existen jerarquia y level; resuelva manualmente en MongoDB")
		return
	}
	from := name + ".jerarquia"
	to := name + ".level"
	err = client.Database("admin").RunCommand(ctx, bson.D{
		{Key: "renameCollection", Value: from},
		{Key: "to", Value: to},
		{Key: "dropTarget", Value: false},
	}).Err()
	if err != nil {
		log.Printf("migración jerarquia: renameCollection falló (%v). Renombre jerarquia → level en Atlas", err)
		return
	}
	log.Printf("migración jerarquia: colección renombrada a %s", to)
}

/* migrateDegreeStudyPlanIDToNumero reemplaza studyplanid ObjectId por numeroresolucion del plan y elimina resolucionid. */
func migrateDegreeStudyPlanIDToNumero(client *mongo.Client) {
	name := strings.TrimSpace(dbName)
	if name == "" {
		name = "san_martin"
	}
	ctx, cancel := context.WithTimeout(context.Background(), 45*time.Second)
	defer cancel()
	db := client.Database(name)
	coll := db.Collection("degree")
	spColl := db.Collection("studyplan")

	cur, err := coll.Find(ctx, bson.M{})
	if err != nil {
		log.Printf("migración degree studyplanid: Find: %v", err)
		return
	}
	defer cur.Close(ctx)

	var updated int64
	for cur.Next(ctx) {
		var doc bson.M
		if err := cur.Decode(&doc); err != nil {
			continue
		}
		idVal, ok := doc["_id"].(primitive.ObjectID)
		if !ok {
			continue
		}
		spv, has := doc["studyplanid"]
		if !has {
			continue
		}
		oid, isOID := spv.(primitive.ObjectID)
		if !isOID {
			continue
		}
		if oid.IsZero() {
			_, _ = coll.UpdateOne(ctx, bson.M{"_id": idVal}, bson.M{"$unset": bson.M{"studyplanid": ""}})
			continue
		}
		var sp models.StudyPlan
		if err := spColl.FindOne(ctx, bson.M{"_id": oid}).Decode(&sp); err != nil {
			log.Printf("migración degree studyplanid: plan %s no encontrado para carrera %s", oid.Hex(), idVal.Hex())
			_, _ = coll.UpdateOne(ctx, bson.M{"_id": idVal}, bson.M{"$unset": bson.M{"studyplanid": ""}})
			continue
		}
		nr := strings.TrimSpace(sp.NumeroResolucion)
		if nr == "" {
			_, _ = coll.UpdateOne(ctx, bson.M{"_id": idVal}, bson.M{"$unset": bson.M{"studyplanid": ""}})
			continue
		}
		if _, err := coll.UpdateOne(ctx, bson.M{"_id": idVal}, bson.M{"$set": bson.M{"studyplanid": nr}}); err == nil {
			updated++
		}
	}
	if err := cur.Err(); err != nil {
		log.Printf("migración degree studyplanid: cursor: %v", err)
	}
	if updated > 0 {
		log.Printf("migración degree studyplanid: carreras actualizadas a número de resolución: %d", updated)
	}
	ur, err := coll.UpdateMany(ctx, bson.M{}, bson.M{"$unset": bson.M{"resolucionid": ""}})
	if err != nil {
		log.Printf("migración degree: quitar resolucionid: %v", err)
		return
	}
	if ur.ModifiedCount > 0 {
		log.Printf("migración degree: campo resolucionid eliminado de documentos: %d", ur.ModifiedCount)
	}
}

/* cleanupDegreeNilStudyPlanIDs quita studyplanid = ObjectId cero (referencia inválida). */
func cleanupDegreeNilStudyPlanIDs(client *mongo.Client) {
	name := strings.TrimSpace(dbName)
	if name == "" {
		name = "san_martin"
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	coll := client.Database(name).Collection("degree")
	ur, err := coll.UpdateMany(ctx, bson.M{"studyplanid": primitive.NilObjectID}, bson.M{"$unset": bson.M{"studyplanid": ""}})
	if err != nil {
		log.Printf("cleanup degree studyplanid nulo: %v", err)
		return
	}
	if ur.ModifiedCount > 0 {
		log.Printf("cleanup degree studyplanid nulo: documentos corregidos: %d", ur.ModifiedCount)
	}
}

/* ConnectDB release a connection to the DB */
func ConnectDB() *mongo.Client {

	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err.Error())
		return client
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err.Error())
		return client
	}

	migratePursueTypeToAttendanceMode(client)
	migrateJerarquiaToLevel(client)
	migrateDegreeStudyPlanIDToNumero(client)
	cleanupDegreeNilStudyPlanIDs(client)

	log.Println("Connection DB successfully")
	return client
}

/* CheckConnection check from anywhere the connection live */
func CheckConnection() int {
	err := MongoConnection.Ping(context.TODO(), nil)
	if err != nil {
		return 0
	}
	return 1
}


