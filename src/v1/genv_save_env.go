package genv

import (
	"fmt"
	"os"
)

var (
	genvDir  string // The directory in which the env file is located
	genvPath string // The path for the env file
)

func Save() error {
	if genvDir == "" || genvPath == "" {
		return fmt.Errorf("can't save because genvDir or genvPath is empty")
	}

	if EnvVariables == nil {
		return nil // Nothing to save
	}

	file, err := os.Create(genvPath)
	if err != nil {
		return err
	}

	defer file.Close()

	// Gathering variables
	var lines []string
	for key, value := range EnvVariables {
		if value != "" {
			line := fmt.Sprintf("%s=%s\n", key, value)
			lines = append(lines, line)
		}
	}

	// Writing to file
	for _, line := range lines {
		_, err := file.WriteString(line)
		if err != nil {
			return fmt.Errorf("something went wrong with writing to file, %s", err)
		}
	}

	return nil // Config saved successfully
}
