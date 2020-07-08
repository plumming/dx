package util

import (
	"os"
	"path/filepath"

	"github.com/jenkins-x/jx-logging/pkg/log"
)

func ChillyBinaryLocation() (string, error) {
	return chillyBinaryLocation(os.Executable)
}

func chillyBinaryLocation(osExecutable func() (string, error)) (string, error) {
	processBinary, err := osExecutable()
	if err != nil {
		log.Logger().Debugf("processBinary error %s", err)
		return processBinary, err
	}
	log.Logger().Debugf("processBinary %s", processBinary)
	// make it absolute
	processBinary, err = filepath.Abs(processBinary)
	if err != nil {
		log.Logger().Debugf("processBinary error %s", err)
		return processBinary, err
	}
	log.Logger().Debugf("processBinary %s", processBinary)

	// if the process was started form a symlink go and get the absolute location.
	processBinary, err = filepath.EvalSymlinks(processBinary)
	if err != nil {
		log.Logger().Debugf("processBinary error %s", err)
		return processBinary, err
	}

	log.Logger().Debugf("processBinary %s", processBinary)
	path := filepath.Dir(processBinary)
	log.Logger().Debugf("dir from '%s' is '%s'", processBinary, path)
	return path, nil
}
