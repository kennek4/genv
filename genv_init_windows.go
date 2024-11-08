package genv

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"syscall"
)

/*
Initializes genv with the given app name and file type in the given directory.
If no directory is given, the default directory genv will use will be located in
$HOME/.{appName} or %USERPROFILE%\.{appName} on Linux and Windows respectively.
*/
func Init(appName string, dir ...string) error {

	var baseDir string
	appName = strings.ToTitle(appName)

	switch len(dir) {
	case 0: // Use default dir
		home, err := os.UserHomeDir()
		if err != nil {
			return ErrGenvFilePathError
		}
		baseDir = home

	case 1: // Use given dir
		path, err := os.Stat(dir[0])
		if err != nil || !path.IsDir() || os.IsNotExist(err) {
			return ErrGenvInvalidDir
		}

		baseDir = dir[0]

	default: // Too many arguments were given
		return fmt.Errorf("too many arguments for dir were given for genv.Init(), only 1 is required")
	}

	// Creates the path with the directory of the app's genv
	genvDir = filepath.Join(baseDir, "."+appName)

	err := os.Mkdir(genvDir, 0777)
	if err != nil {
		return ErrGenvMkValidDir
	}

	// Make genvDir hidden on Windows
	// Currently UNTESTED on Windows machines
	genvDirPtr, err := syscall.UTF16PtrFromString(genvDir)
	if err != nil {
		return err
	}

	err = syscall.SetFileAttributes(genvDirPtr, syscall.FILE_ATTRIBUTE_HIDDEN)
	if err != nil {
		return err
	}

	genvFile := fmt.Sprintf(".%s.env", appName)
	genvPath = filepath.Join(genvDir, genvFile)

	err = os.WriteFile(genvPath, nil, 0777)
	if err != nil {
		return ErrGenvFilePathError
	}

	return nil // genv created successfully
}
