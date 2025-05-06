package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/Beluga-Whale/management-api/internal/models"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB


var TestDB *gorm.DB

func LoadEnv() {
	env := os.Getenv("APP_ENV")

	if env == "production" {
		fmt.Println("✅ Running in production mode: using ENV variables only")
		return
	}

	if env == "" {
		env = "development"
	}

	envFileMap := map[string]string{
		"development":    ".env",
		"test":           ".env.test",
		"test.localhost": ".env.test.localhost",
		"production":     ".env.production",
	}

	envFile, ok := envFileMap[env]
	if !ok {
		log.Fatalf("❌ Invalid APP_ENV: %s", env)
	}

	// ✅ ใช้ runtime.Caller เพื่อให้ไม่หลุด directory เวลา go test
	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		log.Fatal("❌ Failed to get current file path")
	}

	// currentFile → /path/to/project/server/config/config.go
	serverDir := filepath.Join(filepath.Dir(currentFile), "..") // เดินขึ้นจาก /config → /server
	envPath := filepath.Join(serverDir, envFile)

	fmt.Println("🔧 Loading env from:", envPath)

	if err := godotenv.Load(envPath); err != nil {
		log.Fatalf("❌ Failed to load env: %v", err)
	}
}


func ConnectDB() {
	var err error

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
	os.Getenv("HOST"),
	os.Getenv("USER_NAME"),
	os.Getenv("PASSWORD"),
	os.Getenv("DATABASE_NAME"),
	os.Getenv("PORT"),
	os.Getenv("SSL_MODE"),
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
