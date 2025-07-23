package utils

import (
	"encoding/json"
	"strings"
	"testing"
)

// Sample swagger JSON for testing
const sampleSwaggerJSON = `{
  "swagger": "2.0",
  "info": {
    "title": "Test API",
    "version": "1.0"
  },
  "host": "localhost:3000",
  "basePath": "/",
  "paths": {
    "/v1/examples": {
      "get": {
        "tags": ["Example"],
        "summary": "Get examples",
        "responses": {
          "200": {
            "description": "OK"
          }
        }
      },
      "post": {
        "tags": ["Example"],
        "summary": "Create example",
        "responses": {
          "201": {
            "description": "Created"
          }
        }
      }
    },
    "/v1/profile": {
      "get": {
        "tags": ["Access"],
        "summary": "Get profile",
        "responses": {
          "200": {
            "description": "OK"
          }
        }
      }
    },
    "/v1/permissions": {
      "get": {
        "tags": ["Permission"],
        "summary": "Get permissions",
        "responses": {
          "200": {
            "description": "OK"
          }
        }
      }
    },
    "/v1/audit-logs": {
      "get": {
        "tags": ["Audit"],
        "summary": "Get audit logs",
        "responses": {
          "200": {
            "description": "OK"
          }
        }
      }
    }
  }
}`

func TestFilterSwagger(t *testing.T) {
	data := []byte(sampleSwaggerJSON)

	// Test filtering Access and Permission tags
	filteredData, err := FilterSwagger(data, "Access,Permission", false)
	if err != nil {
		t.Fatalf("FilterSwagger failed: %v", err)
	}

	// Parse filtered result
	var result map[string]interface{}
	if err := json.Unmarshal(filteredData, &result); err != nil {
		t.Fatalf("Failed to parse filtered result: %v", err)
	}

	// Check that paths exist
	paths, ok := result["paths"].(map[string]interface{})
	if !ok {
		t.Fatal("Paths not found in filtered result")
	}

	// Should keep /v1/examples and /v1/audit-logs
	if _, exists := paths["/v1/examples"]; !exists {
		t.Error("Expected /v1/examples to be kept")
	}

	if _, exists := paths["/v1/audit-logs"]; !exists {
		t.Error("Expected /v1/audit-logs to be kept")
	}

	// Should remove /v1/profile and /v1/permissions
	if _, exists := paths["/v1/profile"]; exists {
		t.Error("Expected /v1/profile to be removed")
	}

	if _, exists := paths["/v1/permissions"]; exists {
		t.Error("Expected /v1/permissions to be removed")
	}
}

func TestFilterSwaggerWithOptions(t *testing.T) {
	data := []byte(sampleSwaggerJSON)

	options := FilterSwaggerOptions{
		RemoveTags: []string{"Example", "Access"},
		Verbose:    false,
	}

	filteredData, err := FilterSwaggerWithOptions(data, options)
	if err != nil {
		t.Fatalf("FilterSwaggerWithOptions failed: %v", err)
	}

	// Parse filtered result
	var result map[string]interface{}
	if err := json.Unmarshal(filteredData, &result); err != nil {
		t.Fatalf("Failed to parse filtered result: %v", err)
	}

	paths, ok := result["paths"].(map[string]interface{})
	if !ok {
		t.Fatal("Paths not found in filtered result")
	}

	// Should remove /v1/examples and /v1/profile
	if _, exists := paths["/v1/examples"]; exists {
		t.Error("Expected /v1/examples to be removed")
	}

	if _, exists := paths["/v1/profile"]; exists {
		t.Error("Expected /v1/profile to be removed")
	}

	// Should keep /v1/permissions and /v1/audit-logs
	if _, exists := paths["/v1/permissions"]; !exists {
		t.Error("Expected /v1/permissions to be kept")
	}

	if _, exists := paths["/v1/audit-logs"]; !exists {
		t.Error("Expected /v1/audit-logs to be kept")
	}
}

func TestGetSwaggerFilterStats(t *testing.T) {
	data := []byte(sampleSwaggerJSON)

	stats, err := GetSwaggerFilterStats(data, "Access,Permission")
	if err != nil {
		t.Fatalf("GetSwaggerFilterStats failed: %v", err)
	}

	// Check stats
	if stats["total_paths"] != 4 {
		t.Errorf("Expected 4 total paths, got %d", stats["total_paths"])
	}

	// We have 5 endpoints: 2 in /v1/examples (GET, POST), 1 in /v1/profile, 1 in /v1/permissions, 1 in /v1/audit-logs
	if stats["total_endpoints"] != 5 {
		t.Errorf("Expected 5 total endpoints, got %d", stats["total_endpoints"])
	}

	if stats["removed_endpoints"] != 2 {
		t.Errorf("Expected 2 removed endpoints, got %d", stats["removed_endpoints"])
	}

	if stats["kept_endpoints"] != 3 {
		t.Errorf("Expected 3 kept endpoints, got %d", stats["kept_endpoints"])
	}
}

func TestValidateSwaggerJSON(t *testing.T) {
	// Test valid swagger JSON
	data := []byte(sampleSwaggerJSON)
	if err := ValidateSwaggerJSON(data); err != nil {
		t.Errorf("Expected valid swagger JSON, got error: %v", err)
	}

	// Test invalid JSON
	invalidJSON := []byte(`{"invalid": json}`)
	if err := ValidateSwaggerJSON(invalidJSON); err == nil {
		t.Error("Expected error for invalid JSON")
	}

	// Test missing required fields
	missingSwagger := []byte(`{"info": {"title": "Test"}, "paths": {}}`)
	if err := ValidateSwaggerJSON(missingSwagger); err == nil {
		t.Error("Expected error for missing swagger field")
	}

	missingInfo := []byte(`{"swagger": "2.0", "paths": {}}`)
	if err := ValidateSwaggerJSON(missingInfo); err == nil {
		t.Error("Expected error for missing info field")
	}

	missingPaths := []byte(`{"swagger": "2.0", "info": {"title": "Test"}}`)
	if err := ValidateSwaggerJSON(missingPaths); err == nil {
		t.Error("Expected error for missing paths field")
	}
}

func TestFilterSwaggerPretty(t *testing.T) {
	data := []byte(sampleSwaggerJSON)

	filteredData, err := FilterSwaggerPretty(data, "Access,Permission", false)
	if err != nil {
		t.Fatalf("FilterSwaggerPretty failed: %v", err)
	}

	// Check that result is valid JSON
	var result map[string]interface{}
	if err := json.Unmarshal(filteredData, &result); err != nil {
		t.Fatalf("Failed to parse pretty filtered result: %v", err)
	}

	// Check that it's properly formatted (contains indentation)
	// Pretty formatted JSON should contain newlines and spaces for indentation
	prettyStr := string(filteredData)
	if !strings.Contains(prettyStr, "\n") || !strings.Contains(prettyStr, "    ") {
		t.Error("Expected pretty formatted JSON to contain newlines and indentation")
	}
}

func BenchmarkFilterSwagger(b *testing.B) {
	data := []byte(sampleSwaggerJSON)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := FilterSwagger(data, "Access,Permission", false)
		if err != nil {
			b.Fatalf("FilterSwagger failed: %v", err)
		}
	}
}