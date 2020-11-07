package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"

	"github.com/erizkiatama/berat/controllers"
	"github.com/erizkiatama/berat/models"

	"github.com/joho/godotenv"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func initDB(dbHost, dbPort, dbUser, dbName, dbPassword string) *gorm.DB {
	DBURL := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", dbHost, dbPort, dbUser, dbName, dbPassword)
	db, err := gorm.Open("postgres", DBURL)
	if err != nil {
		fmt.Println("Cannot connect to the database")
		log.Fatal(err)
	} else {
		fmt.Println("Connected to the database")
	}

	db.AutoMigrate(&models.Weight{})

	return db
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error getting env: %s", err.Error())
	}

	db := initDB(os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_NAME"), os.Getenv("DB_PASSWORD"))

	weightRepo := models.WeightRepository{DB: db}
	weightController := controllers.WeightController{WeightRepo: weightRepo}

	router := mux.NewRouter()
	router.HandleFunc("/", weightController.Index).Methods("GET")
	router.HandleFunc("/weight/new", weightController.New).Methods("GET")
	router.HandleFunc("/weight/insert", weightController.Insert).Methods("POST")
	router.HandleFunc("/weight/{id}", weightController.Detail).Methods("GET")

	fmt.Println("Listening to port 8000")
	log.Fatal(http.ListenAndServe(":8000", router))
}
