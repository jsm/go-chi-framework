package connector

import (
	"fmt"
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" // Gorm uses this internally
)

// ConnectPostgres initializes a Postgres connection
func ConnectPostgres(host string, username string, password string, port string) *gorm.DB {
	var portConfiguration string
	if port != "" {
		portConfiguration = fmt.Sprintf("port=%s", port)
	}

	configuration := fmt.Sprintf("host=%s %s user=%s dbname=users sslmode=disable password=%s",
		host,
		portConfiguration,
		username,
		password,
	)

	db, err := gorm.Open("postgres", configuration)
	if err != nil {
		log.Println(configuration)
		log.Println("Failed to connect to Postgres")
		log.Fatalln(err)
	}

	// Disable Logging since we handle all errors
	db.LogMode(false)

	// Disabling automatic inclusion of software-generated timestamps
	db.Callback().Create().Remove("gorm:update_time_stamp")
	db.Callback().Update().Remove("gorm:update_time_stamp")
	db.Set("gorm:save_associations", false)

	return db
}
