package utils

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// GenerateUUIDv7 generates a time-sortable UUIDv7
// UUIDv7 format: timestamp (48 bits) + version (4 bits) + random (12 bits) + variant (2 bits) + random (62 bits)
func GenerateUUIDv7() string {
	// Get current timestamp in milliseconds
	timestamp := time.Now().UnixMilli()
	
	// Create a new UUID
	var u [16]byte
	
	// Set timestamp (48 bits = 6 bytes)
	binary.BigEndian.PutUint64(u[0:8], uint64(timestamp)<<16)
	
	// Fill remaining bytes with random data
	rand.Read(u[6:])
	
	// Set version (4 bits) - version 7
	u[6] = (u[6] & 0x0f) | 0x70
	
	// Set variant (2 bits) - RFC 4122 variant
	u[8] = (u[8] & 0x3f) | 0x80
	
	// Format as UUID string
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
		binary.BigEndian.Uint32(u[0:4]),
		binary.BigEndian.Uint16(u[4:6]),
		binary.BigEndian.Uint16(u[6:8]),
		binary.BigEndian.Uint16(u[8:10]),
		u[10:16])
}

// GenerateUUID generates a standard UUID v4 (fallback)
func GenerateUUID() string {
	return uuid.NewString()
}