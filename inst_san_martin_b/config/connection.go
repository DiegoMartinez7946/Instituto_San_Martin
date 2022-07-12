package config

import (
	"context"
	"log"
	"strings"

	"github.com/benjacifre10/san_martin_b/utils"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

/* MongoConnection is a export variable to access anywhere */
var MongoConnection = ConnectDB()

var dbName = utils.GoDotEnvValue("DB_NAME")
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


