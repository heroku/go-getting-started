package inits

import (
	"fmt"
	"log"

	// "os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConDB() {
	var err error

	//dsn := os.Getenv("DB_URL")
	// dsn := os.Getenv("DB_URL")
	dsn := "host=floppy.db.elephantsql.com port=5432 dbname=pvnsgkmb user=pvnsgkmb password=C8Mhor6qxbeWJAeDu-yrIxR6_p4qCKnF sslmode=prefer connect_timeout=10"
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Failed to con DB")
	} else {
		fmt.Println("Con DB Done")
	}

	//    db.AutoMigrate(&models.Album{})
	//    db.AutoMigrate(&models.Person{})

}
