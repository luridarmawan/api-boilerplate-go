package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
)

// SwaggerDoc represents the swagger JSON structure
type SwaggerDoc struct {
	Swagger     string                 `json:"swagger"`
	Info        map[string]interface{} `json:"info"`
	Host        string                 `json:"host"`
	BasePath    string                 `json:"basePath"`
	Paths       map[string]interface{} `json:"paths"`
	Definitions map[string]interface{} `json:"definitions"`
	Security    []interface{}          `json:"securityDefinitions,omitempty"`
}

// PathOperation represents an HTTP operation in swagger
type PathOperation struct {
	Tags        []string               `json:"tags,omitempty"`
	Summary     string                 `json:"summary,omitempty"`
	Description string                 `json:"description,omitempty"`
	Parameters  []interface{}          `json:"parameters,omitempty"`
	Responses   map[string]interface{} `json:"responses,omitempty"`
	Security    []interface{}          `json:"security,omitempty"`
}

func main() {
	var (
		inputFile  = flag.String("input", "docs/swagger.json", "Input swagger JSON file")
		outputFile = flag.String("output", "", "Output swagger JSON file (default: overwrite input)")
		tagsToRemove = flag.String("remove-tags", "Access,Example,Permission", "Comma-separated list of tags to remove")
		verbose    = flag.Bool("verbose", false, "Enable verbose logging")
	)
	flag.Parse()

	// Set output file to input file if not specified
	if *outputFile == "" {
		*outputFile = *inputFile
	}

	// Parse tags to remove
	tags := strings.Split(*tagsToRemove, ",")
	for i, tag := range tags {
		tags[i] = strings.TrimSpace(tag)
	}

	if *verbose {
		fmt.Printf("Input file: %s\n", *inputFile)
		fmt.Printf("Output file: %s\n", *outputFile)
		fmt.Printf("Tags to remove: %v\n", tags)
	}

	// Read swagger JSON file
	data, err := ioutil.ReadFile(*inputFile)
	if err != nil {
		log.Fatalf("Error reading file %s: %v", *inputFile, err)
	}

	// Parse JSON
	var swaggerDoc map[string]interface{}
	if err := json.Unmarshal(data, &swaggerDoc); err != nil {
		log.Fatalf("Error parsing JSON: %v", err)
	}

	// Filter paths
	filteredPaths := filterPaths(swaggerDoc["paths"], tags, *verbose)
	swaggerDoc["paths"] = filteredPaths

	// Convert back to JSON
	filteredData, err := json.MarshalIndent(swaggerDoc, "", "    ")
	if err != nil {
		log.Fatalf("Error marshaling JSON: %v", err)
	}

	// Write to output file
	if err := ioutil.WriteFile(*outputFile, filteredData, 0644); err != nil {
		log.Fatalf("Error writing file %s: %v", *outputFile, err)
	}

	fmt.Printf("✅ Successfully filtered swagger.json\n")
	fmt.Printf("📁 Output: %s\n", *outputFile)
}

func filterPaths(paths interface{}, tagsToRemove []string, verbose bool) map[string]interface{} {
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
			if shouldRemoveOperation(operationMap, tagsToRemove) {
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

	fmt.Printf("📊 Summary: %d endpoints removed, %d endpoints kept\n", removedCount, keptCount)
	return filteredPaths
}

func shouldRemoveOperation(operation map[string]interface{}, tagsToRemove []string) bool {
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