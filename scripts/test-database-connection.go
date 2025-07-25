package main

import (
	"fmt"
	"log"

	"apiserver/configs"
	"apiserver/internal/database"
)

func main() {
	fmt.Println("🔍 Testing Database Connection")
	fmt.Println("==============================")

	// Load configuration
	config := configs.LoadConfig()

	fmt.Printf("📋 Database Configuration:\n")
	fmt.Printf("   Host: %s\n", config.DBHost)
	fmt.Printf("   Port: %s\n", config.DBPort)
	fmt.Printf("   User: %s\n", config.DBUser)
	fmt.Printf("   Database: %s\n", config.DBName)
	fmt.Printf("   SSL Mode: %s\n", config.DBSSLMode)

	// Try to connect to database
	fmt.Println("\n🔄 Attempting to connect...")
	
	database.InitDatabase(config)
	db := database.GetDB()

	// Test the connection
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("❌ Failed to get database instance: %v", err)
	}

	if err := sqlDB.Ping(); err != nil {
		log.Fatalf("❌ Failed to ping database: %v", err)
	}

	fmt.Println("✅ Database connection successful!")
	
	// Check if configurations table exists
	var tableExists bool
	if err := db.Raw("SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_name = 'configurations')").Scan(&tableExists).Error; err != nil {
		fmt.Printf("⚠️  Could not check configurations table: %v\n", err)
	} else if tableExists {
		fmt.Println("✅ Configurations table exists")
	} else {
		fmt.Println("ℹ️  Configurations table does not exist yet (will be created on first run)")
	}

	fmt.Println("\n🎉 Database is ready for configuration module!")
}