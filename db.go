package main

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var DB *gorm.DB

func InitDB() {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		log.Fatal("❌ Failed to connect database: ", err)
	}

	DB = db
	fmt.Println("✅ Database connected")

		err = DB.AutoMigrate(
		&OPD{},
		&JenisPelayanan{},
		&UserOPD{},
		&UserPemda{},
		&FormPemohon{},
		&FormPengajuan{},
	)
	if err != nil {
		log.Fatal("❌ Migration failed: ", err)
	}
	fmt.Println("✅ AutoMigration finished")
}