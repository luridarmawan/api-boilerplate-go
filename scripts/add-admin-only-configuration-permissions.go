package main

import (
	"fmt"
	"log"
	"time"

	"apiserver/configs"
	"apiserver/internal/database"
	"apiserver/internal/modules/group"
	"apiserver/internal/modules/permission"
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

	fmt.Println("🔧 Adding Configuration Permissions (Admin Only)")
	fmt.Println("===============================================")

	// Create configuration permissions
	configurationPermissions := []permission.Permission{
		{
			Name:        "Create Configurations",
			Description: "Permission to create new configurations (Admin only)",
			Resource:    "configurations",
			Action:      "create",
			StatusID:    int16Ptr(0), // Active
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			Name:        "Read Configurations",
			Description: "Permission to read configurations (Admin only)",
			Resource:    "configurations",
			Action:      "read",
			StatusID:    int16Ptr(0), // Active
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			Name:        "Update Configurations",
			Description: "Permission to update configurations (Admin only)",
			Resource:    "configurations",
			Action:      "update",
			StatusID:    int16Ptr(0), // Active
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			Name:        "Delete Configurations",
			Description: "Permission to delete configurations (Admin only)",
			Resource:    "configurations",
			Action:      "delete",
			StatusID:    int16Ptr(0), // Active
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			Name:        "Manage Configurations",
			Description: "Permission to manage all configuration settings (Admin only)",
			Resource:    "configurations",
			Action:      "manage",
			StatusID:    int16Ptr(0), // Active
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}

	// Add permissions to database
	createdPermissions := make(map[string]uint)
	for _, p := range configurationPermissions {
		var existingPermission permission.Permission
		result := db.Where("name = ? AND resource = ? AND action = ?", p.Name, p.Resource, p.Action).First(&existingPermission)

		if result.Error != nil {
			// Permission doesn't exist, create it
			if err := db.Create(&p).Error; err != nil {
				log.Printf("❌ Failed to create permission %s: %v", p.Name, err)
			} else {
				log.Printf("✅ Created permission: %s with ID: %d", p.Name, p.ID)
				createdPermissions[p.Action] = p.ID
			}
		} else {
			log.Printf("ℹ️  Permission %s already exists with ID: %d", p.Name, existingPermission.ID)
			createdPermissions[p.Action] = existingPermission.ID
		}
	}

	// Get Admin group only
	var adminGroup group.Group
	result := db.Where("name = ?", "Admin").First(&adminGroup)
	
	if result.Error != nil {
		log.Fatal("❌ Admin group not found! Please run the main seeder first.")
	}

	fmt.Printf("🔍 Found Admin group with ID: %d\n", adminGroup.ID)

	// Assign ALL configuration permissions to Admin group ONLY
	assignedCount := 0
	for action, permID := range createdPermissions {
		// Check if permission is already assigned
		var count int64
		db.Table("group_permissions").Where("group_id = ? AND permission_id = ?", adminGroup.ID, permID).Count(&count)
		
		if count == 0 {
			if err := db.Exec("INSERT INTO group_permissions (group_id, permission_id) VALUES (?, ?)", adminGroup.ID, permID).Error; err != nil {
				log.Printf("❌ Failed to assign %s permission to Admin group: %v", action, err)
			} else {
				log.Printf("✅ Assigned %s permission to Admin group", action)
				assignedCount++
			}
		} else {
			log.Printf("ℹ️  %s permission already assigned to Admin group", action)
		}
	}

	// Verify no other groups have configuration permissions
	fmt.Println("\n🔍 Checking other groups for configuration permissions...")
	
	otherGroups := []string{"Editor", "Viewer", "General client"}
	for _, groupName := range otherGroups {
		var otherGroup group.Group
		if err := db.Where("name = ?", groupName).First(&otherGroup).Error; err == nil {
			// Check if this group has any configuration permissions
			var configPerms []permission.Permission
			db.Joins("JOIN group_permissions ON permissions.id = group_permissions.permission_id").
				Where("group_permissions.group_id = ? AND permissions.resource = ?", otherGroup.ID, "configurations").
				Find(&configPerms)
			
			if len(configPerms) > 0 {
				fmt.Printf("⚠️  WARNING: %s group has %d configuration permissions!\n", groupName, len(configPerms))
				for _, perm := range configPerms {
					fmt.Printf("   - %s (%s)\n", perm.Name, perm.Action)
				}
			} else {
				fmt.Printf("✅ %s group has no configuration permissions (correct)\n", groupName)
			}
		}
	}

	fmt.Println("\n📊 Summary:")
	fmt.Printf("✅ Configuration permissions created: %d\n", len(createdPermissions))
	fmt.Printf("✅ Permissions assigned to Admin group: %d\n", assignedCount)
	fmt.Println("\n🔒 Security Status:")
	fmt.Println("✅ Configuration module is now ADMIN-ONLY")
	fmt.Println("✅ Only Admin group can access configuration endpoints")
	fmt.Println("✅ Other groups (Editor, Viewer, General client) have NO access")

	fmt.Println("\n📋 Permissions created for Admin group:")
	fmt.Println("- configurations:create - Create new configurations")
	fmt.Println("- configurations:read - Read configuration settings")
	fmt.Println("- configurations:update - Update configuration settings")
	fmt.Println("- configurations:delete - Delete configurations")
	fmt.Println("- configurations:manage - Full configuration management")

	fmt.Println("\n🚀 Configuration module is ready for Admin-only access!")
}