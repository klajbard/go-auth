package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"gopkg.in/mgo.v2"
)

var DB *mgo.Database
var Users *mgo.Collection

func init() {

	err := godotenv.Load(".ENV")
	if err != nil {
		log.Println("Error loading .ENV file")
	}
	log.Println("Loaded env variables")

	s, err := mgo.Dial(os.Getenv("MONGO_URL"))
	if err != nil {
		panic(err)
	}

	if err = s.Ping(); err != nil {
		panic(err)
	}

	DB = s.DB("hassio")
	Users = DB.C("users")

	log.Println("MongoDB connected")
}
