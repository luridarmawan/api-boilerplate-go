package main

import (
	"fmt"
	"log"
	"time"

	"apiserver/configs"
	"apiserver/internal/database"
	"apiserver/internal/modules/configuration"
)

// Helper function to create int16 pointer
func int16Ptr(v int16) *int16 {
	return &v
}

func main() {
	// Load configuration
	config := configs.LoadConfig()

	// Initialize database
	database.InitDatabase(config)
	db := database.GetDB()

	fmt.Println("🌱 Seeding Sample Configuration Data")
	fmt.Println("====================================")

	// Sample configurations
	sampleConfigurations := []configuration.Configuration{
		{
			Key:         "api.rate_limit.default",
			Value:       "120",
			Description: "Default rate limit per minute for API requests",
			StatusID:    int16Ptr(0), // Active
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			Key:         "api.rate_limit.admin",
			Value:       "1000",
			Description: "Rate limit per minute for admin users",
			StatusID:    int16Ptr(0), // Active
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			Key:         "security.jwt.expiry",
			Value:       "24h",
			Description: "JWT token expiry duration",
			StatusID:    int16Ptr(0), // Active
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			Key:         "security.password.min_length",
			Value:       "8",
			Description: "Minimum password length requirement",
			StatusID:    int16Ptr(0), // Active
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			Key:         "database.connection_pool.max",
			Value:       "100",
			Description: "Maximum database connection pool size",
			StatusID:    int16Ptr(0), // Active
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			Key:         "logging.level",
			Value:       "info",
			Description: "Application logging level (debug, info, warn, error)",
			StatusID:    int16Ptr(0), // Active
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			Key:         "features.swagger.enabled",
			Value:       "true",
			Description: "Enable/disable Swagger documentation",
			StatusID:    int16Ptr(0), // Active
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			Key:         "features.audit.enabled",
			Value:       "true",
			Description: "Enable/disable audit logging",
			StatusID:    int16Ptr(0), // Active
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			Key:         "maintenance.mode",
			Value:       "false",
			Description: "Enable/disable maintenance mode",
			StatusID:    int16Ptr(0), // Active
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}

	// Create configurations
	createdCount := 0
	skippedCount := 0

	for _, config := range sampleConfigurations {
		var existingConfig configuration.Configuration
		result := db.Where("key = ?", config.Key).First(&existingConfig)

		if result.Error != nil {
			// Configuration doesn't exist, create it
			if err := db.Create(&config).Error; err != nil {
				log.Printf("❌ Failed to create configuration %s: %v", config.Key, err)
			} else {
				log.Printf("✅ Created configuration: %s = %s", config.Key, config.Value)
				createdCount++
			}
		} else {
			log.Printf("ℹ️  Configuration %s already exists, skipping...", config.Key)
			skippedCount++
		}
	}

	fmt.Printf("\n📊 Seeding Summary:\n")
	fmt.Printf("✅ Configurations created: %d\n", createdCount)
	fmt.Printf("ℹ️  Configurations skipped: %d\n", skippedCount)
	fmt.Printf("📋 Total configurations: %d\n", len(sampleConfigurations))

	// Display all configurations
	fmt.Println("\n📋 Sample Configurations Created:")
	fmt.Println("=================================")

	var allConfigs []configuration.Configuration
	db.Where("status_id = ?", 0).Order("key").Find(&allConfigs)

	for _, config := range allConfigs {
		fmt.Printf("🔑 %s\n", config.Key)
		fmt.Printf("   💾 Value: %s\n", config.Value)
		fmt.Printf("   📝 Description: %s\n", config.Description)
		fmt.Println()
	}

	fmt.Println("🎉 Sample configuration data seeded successfully!")
	fmt.Println("\n📋 Usage Examples:")
	fmt.Println("==================")
	fmt.Println("# Get all configurations")
	fmt.Println(`curl -X GET "http://localhost:3000/v1/configurations" \`)
	fmt.Println(`  -H "Authorization: Bearer admin-api-key-789"`)
	fmt.Println()
	fmt.Println("# Get configuration by key")
	fmt.Println(`curl -X GET "http://localhost:3000/v1/configurations/key/app.name" \`)
	fmt.Println(`  -H "Authorization: Bearer admin-api-key-789"`)
	fmt.Println()
	fmt.Println("# Create new configuration")
	fmt.Println(`curl -X POST "http://localhost:3000/v1/configurations" \`)
	fmt.Println(`  -H "Authorization: Bearer admin-api-key-789" \`)
	fmt.Println(`  -H "Content-Type: application/json" \`)
	fmt.Println(`  -d '{"key":"custom.setting","value":"custom_value","description":"Custom setting"}'`)
	fmt.Println()
	fmt.Println("# Update configuration")
	fmt.Println(`curl -X PUT "http://localhost:3000/v1/configurations/{id}" \`)
	fmt.Println(`  -H "Authorization: Bearer admin-api-key-789" \`)
	fmt.Println(`  -H "Content-Type: application/json" \`)
	fmt.Println(`  -d '{"value":"new_value","description":"Updated description"}'`)
}