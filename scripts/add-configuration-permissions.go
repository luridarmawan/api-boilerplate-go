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

	// Create configuration permissions
	configurationPermissions := []permission.Permission{
		{
			Name:        "Create Configurations",
			Description: "Permission to create new configurations",
			Resource:    "configurations",
			Action:      "create",
			StatusID:    int16Ptr(0), // Active
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			Name:        "Read Configurations",
			Description: "Permission to read configurations",
			Resource:    "configurations",
			Action:      "read",
			StatusID:    int16Ptr(0), // Active
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			Name:        "Update Configurations",
			Description: "Permission to update configurations",
			Resource:    "configurations",
			Action:      "update",
			StatusID:    int16Ptr(0), // Active
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			Name:        "Delete Configurations",
			Description: "Permission to delete configurations",
			Resource:    "configurations",
			Action:      "delete",
			StatusID:    int16Ptr(0), // Active
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}

	// Add permissions to database
	createdPermissions := make(map[string]uint)
	for _, p := range configurationPermissions {
		var existingPermission permission.Permission
		result := db.Where("name = ?", p.Name).First(&existingPermission)

		if result.Error != nil {
			// Permission doesn't exist, create it
			if err := db.Create(&p).Error; err != nil {
				log.Printf("Failed to create permission %s: %v", p.Name, err)
			} else {
				log.Printf("Created permission: %s with ID: %d", p.Name, p.ID)
				createdPermissions[p.Action] = p.ID
			}
		} else {
			log.Printf("Permission %s already exists with ID: %d", p.Name, existingPermission.ID)
			createdPermissions[p.Action] = existingPermission.ID
		}
	}

	// Get groups
	var adminGroup, editorGroup, viewerGroup, generalClientGroup group.Group
	db.Where("name = ?", "Admin").First(&adminGroup)
	db.Where("name = ?", "Editor").First(&editorGroup)
	db.Where("name = ?", "Viewer").First(&viewerGroup)
	db.Where("name = ?", "General client").First(&generalClientGroup)

	// Assign permissions to groups
	if adminGroup.ID > 0 {
		// Admin gets all permissions
		for _, permID := range createdPermissions {
			var perm permission.Permission
			if err := db.First(&perm, permID).Error; err == nil {
				if err := db.Exec("INSERT INTO group_permissions (group_id, permission_id) VALUES (?, ?) ON CONFLICT DO NOTHING", adminGroup.ID, permID).Error; err != nil {
					log.Printf("Failed to assign permission %d to Admin group: %v", permID, err)
				} else {
					log.Printf("Assigned permission %d to Admin group", permID)
				}
			}
		}
	}

	if editorGroup.ID > 0 {
		// Editor gets create, read, update
		for _, action := range []string{"create", "read", "update"} {
			if permID, exists := createdPermissions[action]; exists {
				if err := db.Exec("INSERT INTO group_permissions (group_id, permission_id) VALUES (?, ?) ON CONFLICT DO NOTHING", editorGroup.ID, permID).Error; err != nil {
					log.Printf("Failed to assign permission %d to Editor group: %v", permID, err)
				} else {
					log.Printf("Assigned permission %d to Editor group", permID)
				}
			}
		}
	}

	if viewerGroup.ID > 0 {
		// Viewer gets read only
		if permID, exists := createdPermissions["read"]; exists {
			if err := db.Exec("INSERT INTO group_permissions (group_id, permission_id) VALUES (?, ?) ON CONFLICT DO NOTHING", viewerGroup.ID, permID).Error; err != nil {
				log.Printf("Failed to assign permission %d to Viewer group: %v", permID, err)
			} else {
				log.Printf("Assigned permission %d to Viewer group", permID)
			}
		}
	}

	if generalClientGroup.ID > 0 {
		// General client gets read only
		if permID, exists := createdPermissions["read"]; exists {
			if err := db.Exec("INSERT INTO group_permissions (group_id, permission_id) VALUES (?, ?) ON CONFLICT DO NOTHING", generalClientGroup.ID, permID).Error; err != nil {
				log.Printf("Failed to assign permission %d to General client group: %v", permID, err)
			} else {
				log.Printf("Assigned permission %d to General client group", permID)
			}
		}
	}

	fmt.Println("Configuration permissions have been added and assigned to groups!")
	fmt.Println("\nPermissions created:")
	fmt.Println("- Create Configurations (Admin, Editor)")
	fmt.Println("- Read Configurations (Admin, Editor, Viewer, General client)")
	fmt.Println("- Update Configurations (Admin, Editor)")
	fmt.Println("- Delete Configurations (Admin only)")
}
