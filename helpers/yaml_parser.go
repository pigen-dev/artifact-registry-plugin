package helpers

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

func YAMLParser(filePath string, output interface{}) error {
	if output == nil {
		return fmt.Errorf("output instance is nil cant parse yaml file")
	}

	// Read the YAML file
	f, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read YAML file: %w", err)
	}

	// Unmarshal the YAML content into the pipeline object
	err = yaml.Unmarshal(f, output)
	if err != nil {
		return fmt.Errorf("failed to unmarshal YAML: %w", err)
	}

	return nil
}