package db

import (
	"log"
	"os"

	m "example.com/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Database is the database connection
var Database *gorm.DB

// Connect to database
func Connect() error {
	var err error
	dbURL := os.Getenv("DB_URL")
	if err != nil {
		log.Fatalln("Check your .env file")
	}

	Database, err = gorm.Open(sqlite.Open(dbURL), &gorm.Config{
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
	})
	if err != nil {
		panic(err)
	}

	// https://gorm.io/docs/migration.html
	Database.AutoMigrate(&m.Dog{})

	initData(Database)

	return nil
}

func initData(db *gorm.DB) {
	//You can insert multiple records too
	var dogs []m.Dog = []m.Dog{
		{Name: "Ricky", Breed: "Chihuahua", Age: m.ToNullInt16(2), IsGoodBoy: false},
		{Name: "Adam", Breed: "Pug", IsGoodBoy: true},
		{Name: "Justin", Breed: "Poodle", Age: m.ToNullInt16(3), IsGoodBoy: false},
	}
	tx := db.Create(&dogs)
	if tx.Error != nil {
		log.Fatalln(tx.Error)
	}
	log.Printf("init: %d records inserted", tx.RowsAffected)
}

// GetTableJSONTags get table name and fields' json tags
func GetTableJSONTags(db *gorm.DB, model interface{}) (string, map[string]string) {
	stmt := &gorm.Statement{DB: db}
	if err := stmt.Parse(&model); err != nil {
		log.Fatalf("Model Parse: %s", err)
		return "", nil
	}
	m := make(map[string]string)
	for _, field := range stmt.Schema.Fields {
		m[field.Name] = field.StructField.Tag.Get("json")
	}
	return stmt.Schema.Table, m
}
