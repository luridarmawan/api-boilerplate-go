package database

import (
	"log"
	"time"

	"apiserver/internal/modules/access"
	"apiserver/internal/modules/example"
	"apiserver/internal/modules/group"
	"apiserver/internal/modules/permission"

	"gorm.io/gorm"
)

// Helper function to create int16 pointer
func int16Ptr(v int16) *int16 {
	return &v
}

// Helper function to create uint pointer
func uintPtr(v uint) *uint {
	return &v
}

func SeedDatabase(db *gorm.DB) {
	log.Println("Starting database seeding...")

	// Seed permissions first
	seedPermissions(db)

	// Seed groups
	seedGroups(db)

	// Seed users
	seedUsers(db)

	// Seed examples
	seedExamples(db)

	log.Println("Database seeding completed!")
}

// Helper function to create time pointer
func timePtr(t time.Time) *time.Time {
	return &t
}

func seedUsers(db *gorm.DB) {
	// Calculate dates for expiration examples
	now := time.Now()
	futureDate := now.AddDate(0, 3, 0)  // 3 months in future
	pastDate := now.AddDate(0, -1, 0)   // 1 month in past

	users := []access.User{
		{
			Name:       "Admin User",
			Email:      "admin@example.com",
			APIKey:     "admin-api-key-789",
			GroupID:    uintPtr(1),
			StatusID:   int16Ptr(0), // Active
			RateLimit:  1000,        // High rate limit for admin
			// No ExpiredDate = never expires
		},
		{
			Name:       "John Doe",
			Email:      "john@example.com",
			APIKey:     "test-api-key-123",
			GroupID:    uintPtr(4),
			StatusID:   int16Ptr(0), // Active
			RateLimit:  120,          // Set to 10 for testing rate limit
			ExpiredDate: timePtr(futureDate), // Expires in 3 months
		},
		{
			Name:       "Jane Smith",
			Email:      "jane@example.com",
			RateLimit:  60,          // Lower rate limit
			ExpiredDate: timePtr(pastDate), // Already expired
			APIKey:   "test-api-key-456",
			GroupID:   uintPtr(4),
			StatusID:  int16Ptr(0), // Active
		},
	}

	for _, u := range users {
		var existingUser access.User
		result := db.Where("email = ?", u.Email).First(&existingUser)
		
		if result.Error == gorm.ErrRecordNotFound {
			// Debug: Print status_id before create
			log.Printf("Creating user %s with status_id: %d", u.Email, u.StatusID)

			if err := db.Select("*").Create(&u).Error; err != nil {
				log.Printf("Failed to create user %s: %v", u.Email, err)
			} else {
				log.Printf("Created user: %s with API key: %s", u.Email, u.APIKey)

				// Debug: Verify status_id after create
				var createdUser access.User
				if err := db.Where("email = ?", u.Email).First(&createdUser).Error; err == nil {
					log.Printf("Verified user %s has status_id: %d", createdUser.Email, createdUser.StatusID)
				}
			}
		} else {
			log.Printf("User %s already exists, updating status_id and rate_limit...", u.Email)
			// Update existing user to ensure status_id is 0 and rate_limit is correct
			updates := map[string]interface{}{
				"status_id":  0,
				"rate_limit": u.RateLimit,
			}
			if err := db.Model(&existingUser).Updates(updates).Error; err != nil {
				log.Printf("Failed to update user %s: %v", u.Email, err)
			} else {
				log.Printf("Updated user %s status_id to 0 and rate_limit to %d", u.Email, u.RateLimit)
			}
		}
	}
}

func seedExamples(db *gorm.DB) {
	examples := []example.Example{
		{
			Name:        "First Example",
			Description: "This is the first example for testing",
			StatusID:    int16Ptr(0), // Active
		},
		{
			Name:        "Second Example",
			Description: "This is the second example with more details",
			StatusID:    int16Ptr(0), // Active
		},
		{
			Name:        "API Testing Example",
			Description: "Example specifically created for API testing purposes",
			StatusID:    int16Ptr(0), // Active
		},
	}

	for _, e := range examples {
		var existingExample example.Example
		result := db.Where("name = ?", e.Name).First(&existingExample)
		
		if result.Error == gorm.ErrRecordNotFound {
			if err := db.Create(&e).Error; err != nil {
				log.Printf("Failed to create example %s: %v", e.Name, err)
			} else {
				log.Printf("Created example: %s", e.Name)
			}
		} else {
			log.Printf("Example %s already exists, skipping...", e.Name)
		}
	}
}

func seedPermissions(db *gorm.DB) {
	permissions := []permission.Permission{
		{
			Name:        "Create Examples",
			Description: "Permission to create new examples",
			Resource:    "examples",
			Action:      "create",
			StatusID:    int16Ptr(0), // Active
		},
		{
			Name:        "Read Examples",
			Description: "Permission to read examples",
			Resource:    "examples",
			Action:      "read",
			StatusID:    int16Ptr(0), // Active
		},
		{
			Name:        "Update Examples",
			Description: "Permission to update examples",
			Resource:    "examples",
			Action:      "update",
			StatusID:    int16Ptr(0), // Active
		},
		{
			Name:        "Delete Examples",
			Description: "Permission to delete examples",
			Resource:    "examples",
			Action:      "delete",
			StatusID:    int16Ptr(0), // Active
		},
		{
			Name:        "Manage Permissions",
			Description: "Permission to manage system permissions",
			Resource:    "permissions",
			Action:      "manage",
			StatusID:    int16Ptr(0), // Active
		},
		{
			Name:        "Manage Groups",
			Description: "Permission to manage user groups",
			Resource:    "groups",
			Action:      "manage",
			StatusID:    int16Ptr(0), // Active
		},
		{
			Name:        "View Profile",
			Description: "Permission to view user profile",
			Resource:    "profile",
			Action:      "read",
			StatusID:    int16Ptr(0), // Active
		},
		{
			Name:        "General Access",
			Description: "General API Access",
			Resource:    "general",
			Action:      "read",
			StatusID:    int16Ptr(0), // Active
		},
		{
			Name:        "Read Audit Logs",
			Description: "Permission to read audit logs",
			Resource:    "audit",
			Action:      "read",
			StatusID:    int16Ptr(0), // Active
		},
		{
			Name:        "Manage Audit Logs",
			Description: "Permission to manage audit logs (cleanup)",
			Resource:    "audit",
			Action:      "manage",
			StatusID:    int16Ptr(0), // Active
		},
		{
			Name:        "Manage Access",
			Description: "Permission to manage API keys and expiration dates",
			Resource:    "access",
			Action:      "manage",
			StatusID:    int16Ptr(0), // Active
		},
		// Configuration permissions (Admin only)
		{
			Name:        "Create Configurations",
			Description: "Permission to create new configurations (Admin only)",
			Resource:    "configurations",
			Action:      "create",
			StatusID:    int16Ptr(0), // Active
		},
		{
			Name:        "Read Configurations",
			Description: "Permission to read configurations (Admin only)",
			Resource:    "configurations",
			Action:      "read",
			StatusID:    int16Ptr(0), // Active
		},
		{
			Name:        "Update Configurations",
			Description: "Permission to update configurations (Admin only)",
			Resource:    "configurations",
			Action:      "update",
			StatusID:    int16Ptr(0), // Active
		},
		{
			Name:        "Delete Configurations",
			Description: "Permission to delete configurations (Admin only)",
			Resource:    "configurations",
			Action:      "delete",
			StatusID:    int16Ptr(0), // Active
		},
		{
			Name:        "Manage Configurations",
			Description: "Permission to manage all configuration settings (Admin only)",
			Resource:    "configurations",
			Action:      "manage",
			StatusID:    int16Ptr(0), // Active
		},
	}

	for _, p := range permissions {
		var existingPermission permission.Permission
		result := db.Where("name = ?", p.Name).First(&existingPermission)

		if result.Error == gorm.ErrRecordNotFound {
			if err := db.Create(&p).Error; err != nil {
				log.Printf("Failed to create permission %s: %v", p.Name, err)
			} else {
				log.Printf("Created permission: %s", p.Name)
			}
		} else {
			log.Printf("Permission %s already exists, skipping...", p.Name)
		}
	}
}

func seedGroups(db *gorm.DB) {
	// First, get all permissions
	var allPermissions []permission.Permission
	db.Find(&allPermissions)

	// Create permission maps for easier assignment
	permissionMap := make(map[string]uint)
	for _, p := range allPermissions {
		permissionMap[p.Name] = p.ID
	}

	groups := []struct {
		Group       group.Group
		Permissions []string
	}{
		{
			Group: group.Group{
				Name:        "Admin",
				Description: "Full access to all resources",
				StatusID:    int16Ptr(0), // Active
			},
			Permissions: []string{
				"Create Examples", "Read Examples", "Update Examples", "Delete Examples",
				"Manage Permissions", "Manage Groups", "View Profile",
				"Read Audit Logs", "Manage Audit Logs", "Manage Access",
				"Create Configurations", "Read Configurations", "Update Configurations",
				"Delete Configurations", "Manage Configurations",
			},
		},
		{
			Group: group.Group{
				Name:        "Editor",
				Description: "Can create, read, and update examples",
				StatusID:    int16Ptr(0), // Active
			},
			Permissions: []string{
				"Create Examples", "Read Examples", "Update Examples", "View Profile",
			},
		},
		{
			Group: group.Group{
				Name:        "Viewer",
				Description: "Read-only access to examples",
				StatusID:    int16Ptr(0), // Active
			},
			Permissions: []string{
				"Read Examples", "View Profile",
			},
		},
		{
			Group: group.Group{
				Name:        "General client",
				Description: "General client group",
				StatusID:    int16Ptr(0), // Active
			},
			Permissions: []string{
				"Read Examples", "General Access",
			},
		},
	}

	for _, g := range groups {
		var existingGroup group.Group
		result := db.Where("name = ?", g.Group.Name).First(&existingGroup)

		if result.Error == gorm.ErrRecordNotFound {
			// Create group
			if err := db.Create(&g.Group).Error; err != nil {
				log.Printf("Failed to create group %s: %v", g.Group.Name, err)
				continue
			}

			// Assign permissions to group
			var groupPermissions []permission.Permission
			for _, permName := range g.Permissions {
				if permID, exists := permissionMap[permName]; exists {
					var perm permission.Permission
					if err := db.First(&perm, permID).Error; err == nil {
						groupPermissions = append(groupPermissions, perm)
					}
				}
			}

			if len(groupPermissions) > 0 {
				if err := db.Model(&g.Group).Association("Permissions").Append(groupPermissions); err != nil {
					log.Printf("Failed to assign permissions to group %s: %v", g.Group.Name, err)
				} else {
					log.Printf("Created group: %s with %d permissions", g.Group.Name, len(groupPermissions))
				}
			} else {
				log.Printf("Created group: %s (no permissions assigned)", g.Group.Name)
			}
		} else {
			log.Printf("Group %s already exists, skipping...", g.Group.Name)
		}
	}

	// Update users with groups
	updateUsersWithGroups(db)
}

func updateUsersWithGroups(db *gorm.DB) {
	// Get groups
	var adminGroup, editorGroup, viewerGroup group.Group
	db.Where("name = ?", "Admin").First(&adminGroup)
	db.Where("name = ?", "Editor").First(&editorGroup)
	db.Where("name = ?", "Viewer").First(&viewerGroup)

	// Update users with groups
	userGroupAssignments := map[string]uint{
		"admin@example.com": adminGroup.ID,
		"john@example.com":  editorGroup.ID,
		"jane@example.com":  viewerGroup.ID,
	}

	for email, groupID := range userGroupAssignments {
		if groupID > 0 {
			if err := db.Model(&access.User{}).Where("email = ?", email).Update("group_id", groupID).Error; err != nil {
				log.Printf("Failed to assign group to user %s: %v", email, err)
			} else {
				log.Printf("Assigned group to user: %s", email)
			}
		}
	}
}