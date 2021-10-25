package main

import (
	// "fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}

type AppProcess struct {
	gorm.Model

	AppID  int
	Status string
	User   string
}

func main() {
	db, err := gorm.Open("postgres", "user=postgres password=kush dbname=gorm sslmode=disable")
	CheckError(err)

	// close the databse after the main function finishes
	defer db.Close()
	// Connect the database , initatiate
	database := db.DB()
	// Check if the connections is made or not
	err = database.Ping()
	CheckError(err)

	//Inserting the data onto the database
	db.AutoMigrate(&AppProcess{})
	// Filling the data onto the database
	var person = &AppProcess{AppID: 2342, Status: "pending", User: "kushagra"}
	// updating the data onto the database
	db.Create(person)
}
