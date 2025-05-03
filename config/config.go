package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Beluga-Whale/management-api/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB


var TestDB *gorm.DB

func ConnectDB() {
	var err error

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
	os.Getenv("HOST"),
	os.Getenv("USER_NAME"),
	os.Getenv("PASSWORD"),
	os.Getenv("DATABASE_NAME"),
	os.Getenv("PORT"),
	)

	// New logger for detailed SQL logging
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
		SlowThreshold: time.Second, // Slow SQL threshold
		LogLevel:      logger.Info, // Log level
		Colorful:      true,        // Enable color
		},
	)

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newLogger, // add Logger
	})

	if err != nil {
		log.Fatal("Fail to connect DB : ",err)
	}

	fmt.Println("Connect DB Success!")

	DB.Exec("CREATE TYPE user_role AS ENUM ('admin', 'user');")
	DB.Exec("CREATE TYPE task_status AS ENUM ('active', 'inactive');")
	DB.Exec("CREATE TYPE task_priority AS ENUM ('low', 'medium', 'high');")

	// AutoMigrate จะตรวจสอบและอัปเดตฐานข้อมูล
	err = DB.AutoMigrate(
		&models.Users{},   // ให้ตรวจสอบตาราง Users
		&models.Tasks{},   // ให้ตรวจสอบตาราง Tasks
	)
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

}

func ConnectTestDB() {
	
	
	var err error

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("HOST"),
		os.Getenv("USER_NAME"),
		os.Getenv("PASSWORD"),
		os.Getenv("DATABASE_NAME"),
		os.Getenv("PORT"),
	)

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second, // Slow SQL threshold
			LogLevel:      logger.Info, // Log level
			Colorful:      true,        // Enable color
		},
	)

	// เชื่อมต่อกับ PostgreSQL สำหรับการทดสอบ
	TestDB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newLogger, // เพิ่ม logger
	})

	if err != nil {
		log.Fatal("Fail to connect to test DB: ", err)
	}

	fmt.Println("Connected to Test DB Successfully!")

	TestDB.Exec("CREATE TYPE user_role AS ENUM ('admin', 'user');")
	TestDB.Exec("CREATE TYPE task_status AS ENUM ('active', 'inactive');")
	TestDB.Exec("CREATE TYPE task_priority AS ENUM ('low', 'medium', 'high');")

	// ใช้ AutoMigrate เพื่ออัปเดตฐานข้อมูลสำหรับการทดสอบ
	err = TestDB.AutoMigrate(
		&models.Users{},   // ให้ตรวจสอบตาราง Users
		&models.Tasks{},   // ให้ตรวจสอบตาราง Tasks
	)
	if err != nil {
		log.Fatal("Failed to migrate database for test:", err)
	}
}
