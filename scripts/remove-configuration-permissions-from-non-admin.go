package main

import (
	"fmt"
	"log"

	"apiserver/configs"
	"apiserver/internal/database"
	"apiserver/internal/modules/group"
	"apiserver/internal/modules/permission"
)

func main() {
	// Load configuration
	config := configs.LoadConfig()

	// Initialize database
	database.InitDatabase(config)
	db := database.GetDB()

	fmt.Println("🔒 Removing Configuration Permissions from Non-Admin Groups")
	fmt.Println("=========================================================")

	// Get all groups except Admin
	var nonAdminGroups []group.Group
	db.Where("name != ?", "Admin").Find(&nonAdminGroups)

	if len(nonAdminGroups) == 0 {
		fmt.Println("ℹ️  No non-admin groups found")
		return
	}

	// Get all configuration permissions
	var configPermissions []permission.Permission
	db.Where("resource = ?", "configurations").Find(&configPermissions)

	if len(configPermissions) == 0 {
		fmt.Println("ℹ️  No configuration permissions found")
		return
	}

	fmt.Printf("🔍 Found %d configuration permissions\n", len(configPermissions))
	fmt.Printf("🔍 Found %d non-admin groups\n", len(nonAdminGroups))

	totalRemoved := 0

	// Remove configuration permissions from each non-admin group
	for _, grp := range nonAdminGroups {
		fmt.Printf("\n📋 Checking group: %s (ID: %d)\n", grp.Name, grp.ID)
		
		groupRemovedCount := 0
		for _, perm := range configPermissions {
			// Check if this group has this permission
			var count int64
			db.Table("group_permissions").Where("group_id = ? AND permission_id = ?", grp.ID, perm.ID).Count(&count)
			
			if count > 0 {
				// Remove the permission
				result := db.Exec("DELETE FROM group_permissions WHERE group_id = ? AND permission_id = ?", grp.ID, perm.ID)
				if result.Error != nil {
					log.Printf("❌ Failed to remove permission %s from %s group: %v", perm.Name, grp.Name, result.Error)
				} else {
					log.Printf("✅ Removed permission: %s (%s) from %s group", perm.Name, perm.Action, grp.Name)
					groupRemovedCount++
					totalRemoved++
				}
			}
		}
		
		if groupRemovedCount == 0 {
			fmt.Printf("✅ %s group had no configuration permissions (already clean)\n", grp.Name)
		} else {
			fmt.Printf("🗑️  Removed %d configuration permissions from %s group\n", groupRemovedCount, grp.Name)
		}
	}

	// Verify Admin group still has configuration permissions
	fmt.Println("\n🔍 Verifying Admin group permissions...")
	var adminGroup group.Group
	if err := db.Where("name = ?", "Admin").First(&adminGroup).Error; err != nil {
		fmt.Println("⚠️  WARNING: Admin group not found!")
	} else {
		var adminConfigPerms []permission.Permission
		db.Joins("JOIN group_permissions ON permissions.id = group_permissions.permission_id").
			Where("group_permissions.group_id = ? AND permissions.resource = ?", adminGroup.ID, "configurations").
			Find(&adminConfigPerms)
		
		fmt.Printf("✅ Admin group has %d configuration permissions (should be 5)\n", len(adminConfigPerms))
		
		if len(adminConfigPerms) > 0 {
			fmt.Println("   Admin permissions:")
			for _, perm := range adminConfigPerms {
				fmt.Printf("   - %s (%s)\n", perm.Name, perm.Action)
			}
		}
	}

	// Final verification - check all groups
	fmt.Println("\n📊 Final Security Audit:")
	fmt.Println("========================")

	allGroups := []string{"Admin", "Editor", "Viewer", "General client"}
	for _, groupName := range allGroups {
		var grp group.Group
		if err := db.Where("name = ?", groupName).First(&grp).Error; err == nil {
			var configPerms []permission.Permission
			db.Joins("JOIN group_permissions ON permissions.id = group_permissions.permission_id").
				Where("group_permissions.group_id = ? AND permissions.resource = ?", grp.ID, "configurations").
				Find(&configPerms)
			
			if groupName == "Admin" {
				if len(configPerms) > 0 {
					fmt.Printf("✅ %s: %d configuration permissions (CORRECT)\n", groupName, len(configPerms))
				} else {
					fmt.Printf("❌ %s: %d configuration permissions (ERROR - should have permissions!)\n", groupName, len(configPerms))
				}
			} else {
				if len(configPerms) == 0 {
					fmt.Printf("✅ %s: %d configuration permissions (CORRECT - should be 0)\n", groupName, len(configPerms))
				} else {
					fmt.Printf("❌ %s: %d configuration permissions (ERROR - should be 0!)\n", groupName, len(configPerms))
				}
			}
		}
	}

	fmt.Println("\n📋 Summary:")
	fmt.Printf("🗑️  Total permissions removed: %d\n", totalRemoved)
	fmt.Println("🔒 Configuration module is now secured for Admin-only access")
	
	if totalRemoved > 0 {
		fmt.Println("\n⚠️  IMPORTANT:")
		fmt.Println("   - Configuration permissions have been removed from non-admin groups")
		fmt.Println("   - Only Admin group can now access configuration endpoints")
		fmt.Println("   - This change takes effect immediately")
		fmt.Println("   - Non-admin users will get 403 Forbidden for configuration endpoints")
	} else {
		fmt.Println("\n✅ No changes needed - configuration permissions were already properly secured")
	}
}