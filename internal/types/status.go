package types

// Status constants for all entities
const (
	StatusActive   int16 = 0 // Active/Default
	StatusInactive int16 = 1 // Deleted/Inactive
	StatusPending  int16 = 2 // Pending (for future use)
	StatusSuspended int16 = 3 // Suspended (for future use)
)

// Status descriptions
var StatusDescriptions = map[int16]string{
	StatusActive:    "Active",
	StatusInactive:  "Inactive/Deleted",
	StatusPending:   "Pending",
	StatusSuspended: "Suspended",
}

// IsActive checks if the status is active
func IsActive(statusID int16) bool {
	return statusID == StatusActive
}

// IsInactive checks if the status is inactive/deleted
func IsInactive(statusID int16) bool {
	return statusID == StatusInactive
}

// GetStatusDescription returns the description for a status ID
func GetStatusDescription(statusID int16) string {
	if desc, exists := StatusDescriptions[statusID]; exists {
		return desc
	}
	return "Unknown"
}