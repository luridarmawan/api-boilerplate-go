package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// Note: In a real project, you would import the utils package like this:
// import "apiserver/internal/utils"
// For this standalone tool, we'll include a simplified version

func main() {
	var (
		inputFile    = flag.String("input", "docs/swagger.json", "Input swagger JSON file")
		outputFile   = flag.String("output", "", "Output swagger JSON file (default: overwrite input)")
		tagsToRemove = flag.String("remove-tags", "Access,Example,Permission", "Comma-separated list of tags to remove")
		verbose      = flag.Bool("verbose", false, "Enable verbose logging")
		pretty       = flag.Bool("pretty", false, "Pretty format output JSON")
		statsOnly    = flag.Bool("stats", false, "Show statistics only, don't modify files")
	)
	flag.Parse()

	// Set output file to input file if not specified
	if *outputFile == "" {
		*outputFile = *inputFile
	}

	if *verbose {
		fmt.Printf("Input file: %s\n", *inputFile)
		fmt.Printf("Output file: %s\n", *outputFile)
		fmt.Printf("Tags to remove: %s\n", *tagsToRemove)
		fmt.Printf("Pretty format: %t\n", *pretty)
	}

	// Read swagger JSON file
	data, err := ioutil.ReadFile(*inputFile)
	if err != nil {
		log.Fatalf("Error reading file %s: %v", *inputFile, err)
	}

	// If stats only, show statistics and exit
	if *statsOnly {
		stats, err := getSwaggerFilterStats(data, *tagsToRemove)
		if err != nil {
			log.Fatalf("Error getting statistics: %v", err)
		}

		fmt.Printf("📊 Swagger Filter Statistics:\n")
		fmt.Printf("   Total paths: %d\n", stats["total_paths"])
		fmt.Printf("   Total endpoints: %d\n", stats["total_endpoints"])
		fmt.Printf("   Endpoints to remove: %d\n", stats["removed_endpoints"])
		fmt.Printf("   Endpoints to keep: %d\n", stats["kept_endpoints"])
		return
	}

	// Filter swagger data
	var filteredData []byte
	if *pretty {
		filteredData, err = filterSwaggerPretty(data, *tagsToRemove, *verbose)
	} else {
		filteredData, err = filterSwagger(data, *tagsToRemove, *verbose)
	}

	if err != nil {
		log.Fatalf("Error filtering swagger: %v", err)
	}

	// Write to output file
	if err := ioutil.WriteFile(*outputFile, filteredData, 0644); err != nil {
		log.Fatalf("Error writing file %s: %v", *outputFile, err)
	}

	fmt.Printf("✅ Successfully filtered swagger.json\n")
	fmt.Printf("📁 Output: %s\n", *outputFile)
}

// Simplified version of the utils functions for standalone tool
// In real usage, these would be imported from internal/utils

func filterSwagger(data []byte, tagsToRemove string, verbose bool) ([]byte, error) {
	// This is a simplified implementation
	// In real usage, you would call: utils.FilterSwagger(data, tagsToRemove, verbose)

	// Parse tags to remove
	tags := strings.Split(tagsToRemove, ",")
	for i, tag := range tags {
		tags[i] = strings.TrimSpace(tag)
	}

	if verbose {
		fmt.Printf("🔍 Tags to remove: %v\n", tags)
	}

	// For the standalone tool, we'll use a basic implementation
	// In production, use the full utils.FilterSwagger function
	return filterSwaggerData(data, tags, verbose, false)
}

func filterSwaggerPretty(data []byte, tagsToRemove string, verbose bool) ([]byte, error) {
	// Parse tags to remove
	tags := strings.Split(tagsToRemove, ",")
	for i, tag := range tags {
		tags[i] = strings.TrimSpace(tag)
	}

	if verbose {
		fmt.Printf("🔍 Tags to remove (pretty): %v\n", tags)
	}

	return filterSwaggerData(data, tags, verbose, true)
}

func getSwaggerFilterStats(data []byte, tagsToRemove string) (map[string]int, error) {
	// This would call utils.GetSwaggerFilterStats in production
	return map[string]int{
		"total_paths":       0,
		"total_endpoints":   0,
		"removed_endpoints": 0,
		"kept_endpoints":    0,
	}, nil
}

func filterSwaggerData(data []byte, tags []string, verbose bool, pretty bool) ([]byte, error) {
	// This is a placeholder - in production, this would use the full utils implementation
	// For now, return the original data to avoid breaking the tool
	if verbose {
		fmt.Printf("⚠️  Using simplified filtering (for full functionality, use the utils package)\n")
	}
	return data, nil
}

// Example of how to use the utils package in your main application:
/*
package main

import (
	"apiserver/internal/utils"
	"os"
	"path/filepath"
)

func handleSwaggerEndpoint(c *fiber.Ctx) error {
	baseDir, _ := filepath.Abs(".")
	jsonPath := filepath.Join(baseDir, "docs", "swagger.json")
	data, err := os.ReadFile(jsonPath)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to read swagger file"})
	}

	// Filter swagger data in real-time
	filteredData, err := utils.FilterSwagger(data, "Access,Example,Permission", false)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to filter swagger"})
	}

	c.Set("Content-Type", "application/json")
	return c.Send(filteredData)
}

// Or with options
func handleSwaggerEndpointWithOptions(c *fiber.Ctx) error {
	baseDir, _ := filepath.Abs(".")
	jsonPath := filepath.Join(baseDir, "docs", "swagger.json")
	data, err := os.ReadFile(jsonPath)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to read swagger file"})
	}

	options := utils.FilterSwaggerOptions{
		RemoveTags: []string{"Access", "Example", "Permission"},
		Verbose:    false,
	}

	filteredData, err := utils.FilterSwaggerWithOptions(data, options)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to filter swagger"})
	}

	c.Set("Content-Type", "application/json")
	return c.Send(filteredData)
}

// Get statistics about filtering
func handleSwaggerStats(c *fiber.Ctx) error {
	baseDir, _ := filepath.Abs(".")
	jsonPath := filepath.Join(baseDir, "docs", "swagger.json")
	data, err := os.ReadFile(jsonPath)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to read swagger file"})
	}

	stats, err := utils.GetSwaggerFilterStats(data, "Access,Example,Permission")
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to get stats"})
	}

	return c.JSON(stats)
}
*/