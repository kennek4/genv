package genv

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

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
		return err
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
