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
	// กำหนด APP_ENV ถ้าไม่ถูกตั้งไว้
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "development"
	}

	// กำหนดชื่อไฟล์ env ใช้คำสั่ง set APP_ENV=โหมดที่อยากใช้
	envFileMap := map[string]string{
		"development":     ".env",
		"test":            ".env.test",
		"test.localhost":  ".env.test.localhost",
		"production":      ".env.production",
	}
	
	envFile, ok := envFileMap[env]
	if !ok {
		log.Fatalf("❌ Invalid APP_ENV: %s", env)
	}

	// ใช้ runtime.Caller เพื่อหา directory ของไฟล์นี้
	// สมมติว่าไฟล์นี้อยู่ใน "server/config" ดังนั้น project root อยู่สองระดับขึ้นไป
	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		log.Fatal("❌ Cannot get current file info")
	}
	// currentFile ตัวอย่าง: C:\Users\...\Desktop\profile\ManageMent\server\config\env.go
	configDir := filepath.Dir(currentFile)
	projectRoot := filepath.Join(configDir, "..", "..") // เดินขึ้นไปสองระดับ → คาดว่าจะเป็น root ของโปรเจกต์

	// จาก project root ให้ระบุให้ชัดเจนว่าไฟล์ env อยู่ในโฟลเดอร์ server
	fullPath := filepath.Join(projectRoot, "server", envFile)

	fmt.Printf("🔧 Loading env from: %s\n", fullPath)
	err := godotenv.Load(fullPath)
	if err != nil {
		log.Fatalf("❌ Failed to load env: %v", err)
	}
}

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
