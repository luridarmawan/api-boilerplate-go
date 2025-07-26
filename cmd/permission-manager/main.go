package main

import (
	"fmt"
	"os"
	"strings"

	"apiserver/configs"
	"apiserver/internal/database"
	"apiserver/internal/modules/access"
	"apiserver/internal/modules/permission"

	"gorm.io/gorm"
)

const (
	ExitSuccess = 0
	ExitError   = 1
)

type PermissionManager struct {
	db *gorm.DB
}

func NewPermissionManager(db *gorm.DB) *PermissionManager {
	return &PermissionManager{db: db}
}

// validateAccessID checks if the access ID exists and returns the user with group
func (pm *PermissionManager) validateAccessID(accessID string) (*access.User, error) {
	var user access.User
	err := pm.db.Preload("Group").Where("id = ? AND status_id = ?", accessID, 0).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("access ID '%s' not found", accessID)
		}
		return nil, fmt.Errorf("failed to query access ID: %v", err)
	}

	if user.Group == nil {
		return nil, fmt.Errorf("access ID '%s' is not assigned to any group", accessID)
	}

	return &user, nil
}

// findPermission finds a permission by resource and action
func (pm *PermissionManager) findPermission(resource, action string) (*permission.Permission, error) {
	var perm permission.Permission
	err := pm.db.Where("resource = ? AND action = ? AND status_id = ?", resource, action, 0).First(&perm).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("permission '%s:%s' not found in database", resource, action)
		}
		return nil, fmt.Errorf("failed to query permission: %v", err)
	}

	return &perm, nil
}

// checkPermissionInGroup checks if a permission already exists in the group
func (pm *PermissionManager) checkPermissionInGroup(groupID uint, permissionID uint) (bool, error) {
	var count int64
	err := pm.db.Table("group_permissions").
		Where("group_id = ? AND permission_id = ?", groupID, permissionID).
		Count(&count).Error

	if err != nil {
		return false, fmt.Errorf("failed to check permission in group: %v", err)
	}

	return count > 0, nil
}

// addPermissionToGroup adds a permission to a group
func (pm *PermissionManager) addPermissionToGroup(groupID uint, permissionID uint) error {
	// Use raw SQL to insert into the many-to-many table
	err := pm.db.Exec("INSERT INTO group_permissions (group_id, permission_id) VALUES (?, ?)", groupID, permissionID).Error
	if err != nil {
		return fmt.Errorf("failed to add permission to group: %v", err)
	}

	return nil
}

// processPermissionRequest handles the main logic of adding permission to group
func (pm *PermissionManager) processPermissionRequest(accessID, resource, action string) error {
	// Step 1: Validate access ID and get user with group
	user, err := pm.validateAccessID(accessID)
	if err != nil {
		return err
	}

	// Step 2: Find the permission
	perm, err := pm.findPermission(resource, action)
	if err != nil {
		return err
	}

	// Step 3: Check if permission already exists in group
	exists, err := pm.checkPermissionInGroup(user.Group.ID, perm.ID)
	if err != nil {
		return err
	}

	if exists {
		fmt.Printf("⚠ Warning: Permission '%s:%s' already exists in group '%s'\n", resource, action, user.Group.Name)
		return nil
	}

	// Step 4: Add permission to group
	err = pm.addPermissionToGroup(user.Group.ID, perm.ID)
	if err != nil {
		return err
	}

	// Step 5: Success message
	fmt.Printf("✓ Permission '%s:%s' successfully added to group '%s' for access ID '%s'\n",
		resource, action, user.Group.Name, accessID)

	return nil
}

func printUsage() {
	fmt.Println("Usage: permission-manager [access_id] [resource] [action]")
	fmt.Println("")
	fmt.Println("Arguments:")
	fmt.Println("  access_id  UUID of the access/user to add permission to")
	fmt.Println("  resource   Resource name (e.g., 'configurations', 'examples')")
	fmt.Println("  action     Action name (e.g., 'create', 'read', 'update', 'delete')")
	fmt.Println("")
	fmt.Println("Example:")
	fmt.Println("  permission-manager 019847a9-4efb-72c1-92fb-2c5eab3335d1 configurations create")
	fmt.Println("")
	fmt.Println("Description:")
	fmt.Println("  This tool adds permissions to the group associated with the specified access ID.")
	fmt.Println("  The permission must already exist in the database.")
}

func validateArguments(args []string) error {
	if len(args) != 4 {
		return fmt.Errorf("invalid number of arguments")
	}

	accessID := strings.TrimSpace(args[1])
	resource := strings.TrimSpace(args[2])
	action := strings.TrimSpace(args[3])

	if accessID == "" {
		return fmt.Errorf("access_id cannot be empty")
	}

	if resource == "" {
		return fmt.Errorf("resource cannot be empty")
	}

	if action == "" {
		return fmt.Errorf("action cannot be empty")
	}

	// Basic UUID format validation (simple check)
	if len(accessID) != 36 || strings.Count(accessID, "-") != 4 {
		return fmt.Errorf("access_id must be a valid UUID format")
	}

	return nil
}

func main() {
	// Check for help flag
	if len(os.Args) > 1 && (os.Args[1] == "-h" || os.Args[1] == "--help") {
		printUsage()
		os.Exit(ExitSuccess)
	}

	// Validate arguments
	if err := validateArguments(os.Args); err != nil {
		fmt.Printf("✗ Error: %s\n\n", err)
		printUsage()
		os.Exit(ExitError)
	}

	accessID := strings.TrimSpace(os.Args[1])
	resource := strings.TrimSpace(os.Args[2])
	action := strings.TrimSpace(os.Args[3])

	// Load configuration
	config := configs.LoadConfig()

	// Initialize database
	database.InitDatabase(config)
	db := database.GetDB()

	// Create permission manager
	pm := NewPermissionManager(db)

	// Process the permission request
	if err := pm.processPermissionRequest(accessID, resource, action); err != nil {
		fmt.Printf("✗ Error: %s\n", err)
		os.Exit(ExitError)
	}

	os.Exit(ExitSuccess)
}
