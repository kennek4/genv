/*
genv is a lightweight package to manage
environment variables in Go programs

Made by Ken Tabanay kentabanay@gmail.com
*/
package genv

import (
	"bufio"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var (
	// An error that is returned when genv fails to write to the env file
	errGenvWrite error = errors.New("failed to write to env file")

	// An error that is returned when genv is given an invalid directory
	errGenvInvalidDir error = errors.New("the provided directory is not valid")

	// An error that is returned when genv fails to create a directory for itself
	errGenvMkValidDir error = errors.New("failed to create the directory")

	// An error that is returned when genv encounters a FilePath error
	errGenvFilePathError error = errors.New("failed to write to env file at given path")
)

var (
	genvDir      string         // The directory in which the env file is located
	genvPath     string         // The path for the env file
	EnvVariables = initEnvMap() // A map containing the environment variables during program execution
)

// Creates an map with string key value pairs
func initEnvMap() map[string]string {
	return make(map[string]string)
}

// Maps the following key with the corresponding value
func addToEnvMap(key string, value string) {
	EnvVariables[key] = value
}

func CreateStringVar(key string, value string) {
	addToEnvMap(key, value)
}

// Converts the given int value to a string before adding to the map using strconv.Itoa()
func CreateIntVar(key string, value int) {
	valueString := strconv.Itoa(value)
	addToEnvMap(key, valueString)
}

// Converts the given float64 value into a string before adding to the map.
// The float64 is converted into a string with no exponents and will use the smallest amount of
// number of digits such that ParseFloat will return the float exactly.
func CreateFloatVar(key string, value float64) {
	valueString := strconv.FormatFloat(value, 'f', -1, 64)
	addToEnvMap(key, valueString)
}

// From the map retrieves the following key and returns its associated value.
func GetVar(key string) (value string) {
	return EnvVariables[key]
}

// Saves all the key value pairs contained in the EnvVariables map and writes them to the corresponding
// file with the app name defined during the genv Init() process.
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
			return errGenvWrite
		}
	}

	return nil // Config saved successfully
}

// Loads the variables found in .appName.env and by default will look to see if a directory named .appName exists in the user's home directory.
// On Linux the default location will be at $HOME. On Windows the default location will be at %USERPROFILE%
// The Load function can also be given a directory explicitly to see if .appName.env exists.
func Load(appName string, dir ...string) error {
	appName = strings.ToTitle(appName)

	var pathToCheck string

	switch len(dir) {
	case 0: // Check default genv location ($HOME/.appName or %USERPROFILE%\.appName)
		pathToCheck, _ = os.UserHomeDir()
	case 1: // Check at the given directory
		if _, err := os.Stat(dir[0]); os.IsNotExist(err) {
			return fmt.Errorf("the the provided path is not valid")
		}
		pathToCheck = dir[0]
	default: // Too many arguments were given
		return fmt.Errorf("too many arguments were given for dir")
	}

	err := filepath.WalkDir(pathToCheck, func(path string, d fs.DirEntry, err error) error {
		if strings.Contains(path, appName) {
			genvDir = filepath.Dir(path)
			return nil
		}

		return nil
	})

	if err != nil {
		return errGenvInvalidDir
	}

	genvPath = filepath.Join(genvDir, "."+appName+".env")
	file, err := os.Open(genvPath)
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.Split(scanner.Text(), "=")
		EnvVariables[line[0]] = line[1] // line[0] = key, line[1] = value
	}

	file.Close()
	return nil

}
