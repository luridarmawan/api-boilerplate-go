package utils

import (
	"encoding/json"
	"fmt"
	"strings"
)

// FilterSwaggerOptions contains options for filtering swagger documentation
type FilterSwaggerOptions struct {
	RemoveTags []string
	Verbose    bool
}

// FilterSwagger filters swagger JSON data by removing endpoints with specified tags
// Parameters:
//   - data: Raw JSON data from swagger.json file
//   - tagsToRemove: Comma-separated string of tags to remove (e.g., "Access,Example,Permission")
//   - verbose: Enable verbose logging (optional)
//
// Returns:
//   - Filtered JSON data as bytes
//   - Error if any
//
// Example usage:
//   data, _ := os.ReadFile("docs/swagger.json")
//   filteredData, err := FilterSwagger(data, "Access,Example,Permission", false)
//   if err != nil {
//       return err
//   }
//   c.Set("Content-Type", "application/json")
//   return c.Send(filteredData)
func FilterSwagger(data []byte, tagsToRemove string, verbose ...bool) ([]byte, error) {
	// Parse verbose option
	isVerbose := false
	if len(verbose) > 0 {
		isVerbose = verbose[0]
	}

	// Parse tags to remove
	tags := strings.Split(tagsToRemove, ",")
	for i, tag := range tags {
		tags[i] = strings.TrimSpace(tag)
	}

	if isVerbose {
		fmt.Printf("🔍 FilterSwagger: Tags to remove: %v\n", tags)
	}

	// Parse JSON
	var swaggerDoc map[string]interface{}
	if err := json.Unmarshal(data, &swaggerDoc); err != nil {
		return nil, fmt.Errorf("error parsing swagger JSON: %v", err)
	}

	// Filter paths
	filteredPaths := filterSwaggerPaths(swaggerDoc["paths"], tags, isVerbose)
	swaggerDoc["paths"] = filteredPaths

	// Convert back to JSON
	filteredData, err := json.Marshal(swaggerDoc)
	if err != nil {
		return nil, fmt.Errorf("error marshaling filtered JSON: %v", err)
	}

	return filteredData, nil
}

// FilterSwaggerWithOptions filters swagger JSON data with detailed options
func FilterSwaggerWithOptions(data []byte, options FilterSwaggerOptions) ([]byte, error) {
	tagsString := strings.Join(options.RemoveTags, ",")
	return FilterSwagger(data, tagsString, options.Verbose)
}

// FilterSwaggerPretty filters swagger JSON and returns pretty-formatted JSON
func FilterSwaggerPretty(data []byte, tagsToRemove string, verbose ...bool) ([]byte, error) {
	// Parse verbose option
	isVerbose := false
	if len(verbose) > 0 {
		isVerbose = verbose[0]
	}

	// Parse tags to remove
	tags := strings.Split(tagsToRemove, ",")
	for i, tag := range tags {
		tags[i] = strings.TrimSpace(tag)
	}

	if isVerbose {
		fmt.Printf("🔍 FilterSwaggerPretty: Tags to remove: %v\n", tags)
	}

	// Parse JSON
	var swaggerDoc map[string]interface{}
	if err := json.Unmarshal(data, &swaggerDoc); err != nil {
		return nil, fmt.Errorf("error parsing swagger JSON: %v", err)
	}

	// Filter paths
	filteredPaths := filterSwaggerPaths(swaggerDoc["paths"], tags, isVerbose)
	swaggerDoc["paths"] = filteredPaths

	// Convert back to pretty JSON
	filteredData, err := json.MarshalIndent(swaggerDoc, "", "    ")
	if err != nil {
		return nil, fmt.Errorf("error marshaling filtered JSON: %v", err)
	}

	return filteredData, nil
}

// GetSwaggerFilterStats returns statistics about filtering operation
func GetSwaggerFilterStats(data []byte, tagsToRemove string) (map[string]int, error) {
	// Parse tags to remove
	tags := strings.Split(tagsToRemove, ",")
	for i, tag := range tags {
		tags[i] = strings.TrimSpace(tag)
	}

	// Parse JSON
	var swaggerDoc map[string]interface{}
	if err := json.Unmarshal(data, &swaggerDoc); err != nil {
		return nil, fmt.Errorf("error parsing swagger JSON: %v", err)
	}

	stats := map[string]int{
		"total_paths":     0,
		"total_endpoints": 0,
		"removed_endpoints": 0,
		"kept_endpoints":   0,
	}

	pathsMap, ok := swaggerDoc["paths"].(map[string]interface{})
	if !ok {
		return stats, nil
	}

	stats["total_paths"] = len(pathsMap)

	for _, pathData := range pathsMap {
		pathMap, ok := pathData.(map[string]interface{})
		if !ok {
			continue
		}

		// Check each HTTP method
		for _, operation := range pathMap {
			operationMap, ok := operation.(map[string]interface{})
			if !ok {
				stats["total_endpoints"]++
				stats["kept_endpoints"]++
				continue
			}

			stats["total_endpoints"]++

			if shouldRemoveSwaggerOperation(operationMap, tags) {
				stats["removed_endpoints"]++
			} else {
				stats["kept_endpoints"]++
			}
		}
	}

	return stats, nil
}

// filterSwaggerPaths filters paths in swagger document
func filterSwaggerPaths(paths interface{}, tagsToRemove []string, verbose bool) map[string]interface{} {
	pathsMap, ok := paths.(map[string]interface{})
	if !ok {
		return make(map[string]interface{})
	}

	filteredPaths := make(map[string]interface{})
	removedCount := 0
	keptCount := 0

	for path, pathData := range pathsMap {
		pathMap, ok := pathData.(map[string]interface{})
		if !ok {
			continue
		}

		filteredPathData := make(map[string]interface{})
		hasValidOperations := false

		// Check each HTTP method (get, post, put, delete, etc.)
		for method, operation := range pathMap {
			operationMap, ok := operation.(map[string]interface{})
			if !ok {
				filteredPathData[method] = operation
				hasValidOperations = true
				continue
			}

			// Check if operation has tags to remove
			if shouldRemoveSwaggerOperation(operationMap, tagsToRemove) {
				if verbose {
					fmt.Printf("🗑️  Removing: %s %s\n", strings.ToUpper(method), path)
				}
				removedCount++
			} else {
				filteredPathData[method] = operation
				hasValidOperations = true
				if verbose {
					fmt.Printf("✅ Keeping: %s %s\n", strings.ToUpper(method), path)
				}
				keptCount++
			}
		}

		// Only include path if it has valid operations
		if hasValidOperations {
			filteredPaths[path] = filteredPathData
		}
	}

	if verbose {
		fmt.Printf("📊 Summary: %d endpoints removed, %d endpoints kept\n", removedCount, keptCount)
	}

	return filteredPaths
}

// shouldRemoveSwaggerOperation checks if an operation should be removed based on tags
func shouldRemoveSwaggerOperation(operation map[string]interface{}, tagsToRemove []string) bool {
	tagsInterface, exists := operation["tags"]
	if !exists {
		return false
	}

	tags, ok := tagsInterface.([]interface{})
	if !ok {
		return false
	}

	// Check if any tag in the operation matches tags to remove
	for _, tagInterface := range tags {
		tag, ok := tagInterface.(string)
		if !ok {
			continue
		}

		for _, removeTag := range tagsToRemove {
			if tag == removeTag {
				return true
			}
		}
	}

	return false
}

// ValidateSwaggerJSON validates if the provided data is valid swagger JSON
func ValidateSwaggerJSON(data []byte) error {
	var swaggerDoc map[string]interface{}
	if err := json.Unmarshal(data, &swaggerDoc); err != nil {
		return fmt.Errorf("invalid JSON format: %v", err)
	}

	// Check for required swagger fields
	if _, exists := swaggerDoc["swagger"]; !exists {
		return fmt.Errorf("missing 'swagger' field in JSON")
	}

	if _, exists := swaggerDoc["info"]; !exists {
		return fmt.Errorf("missing 'info' field in JSON")
	}

	if _, exists := swaggerDoc["paths"]; !exists {
		return fmt.Errorf("missing 'paths' field in JSON")
	}

	return nil
}