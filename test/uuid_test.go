package test

import (
	"apiserver/internal/utils"
	"strings"
	"testing"
	"time"
)

func TestGenerateUUIDv7(t *testing.T) {
	// Generate multiple UUIDs
	uuid1 := utils.GenerateUUIDv7()
	time.Sleep(1 * time.Millisecond) // Small delay to ensure different timestamps
	uuid2 := utils.GenerateUUIDv7()

	// Test basic format
	if len(uuid1) != 36 {
		t.Errorf("Expected UUID length 36, got %d", len(uuid1))
	}

	if len(uuid2) != 36 {
		t.Errorf("Expected UUID length 36, got %d", len(uuid2))
	}

	// Test UUID format (8-4-4-4-12)
	parts1 := strings.Split(uuid1, "-")
	if len(parts1) != 5 {
		t.Errorf("Expected 5 parts separated by hyphens, got %d", len(parts1))
	}

	// Test that UUIDs are different
	if uuid1 == uuid2 {
		t.Error("Generated UUIDs should be different")
	}

	// Test time-sortable property (uuid2 should be greater than uuid1)
	if uuid1 >= uuid2 {
		t.Errorf("UUIDv7 should be time-sortable: %s should be < %s", uuid1, uuid2)
	}

	// Test version (should be 7)
	versionChar := uuid1[14] // Version is at position 14
	if versionChar != '7' {
		t.Errorf("Expected version 7, got %c", versionChar)
	}

	t.Logf("Generated UUIDs: %s, %s", uuid1, uuid2)
}

func TestGenerateUUID(t *testing.T) {
	uuid1 := utils.GenerateUUID()
	uuid2 := utils.GenerateUUID()

	// Test basic format
	if len(uuid1) != 36 {
		t.Errorf("Expected UUID length 36, got %d", len(uuid1))
	}

	// Test that UUIDs are different
	if uuid1 == uuid2 {
		t.Error("Generated UUIDs should be different")
	}

	t.Logf("Generated standard UUIDs: %s, %s", uuid1, uuid2)
}

func TestAuditLogUUIDGeneration(t *testing.T) {
	// Test that audit log UUIDs are time-sortable
	uuid1 := utils.GenerateUUIDv7()
	time.Sleep(1 * time.Millisecond) // Ensure different timestamps
	uuid2 := utils.GenerateUUIDv7()
	time.Sleep(1 * time.Millisecond) // Ensure different timestamps
	uuid3 := utils.GenerateUUIDv7()

	// All should be different
	if uuid1 == uuid2 || uuid2 == uuid3 || uuid1 == uuid3 {
		t.Error("All generated UUIDs should be unique")
	}

	// Should be sortable (lexicographically ordered by time)
	if !(uuid1 < uuid2 && uuid2 < uuid3) {
		t.Errorf("UUIDs should be time-sortable: %s < %s < %s", uuid1, uuid2, uuid3)
	}

	t.Logf("Time-sortable audit log UUIDs: %s, %s, %s", uuid1, uuid2, uuid3)
}