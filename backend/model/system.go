package model

import (
	"github.com/apex/log"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"strings"
)

var executableNameMap = map[string]string{
	// python
	"python":    "Python",
	"python2":   "Python 2",
	"python2.7": "Python 2.7",
	"python3":   "Python 3",
	"python3.5": "Python 3.5",
	"python3.6": "Python 3.6",
	"python3.7": "Python 3.7",
	"python3.8": "Python 3.8",
	// java
	"java": "Java",
	// go
	"go": "Go",
	// node
	"node": "NodeJS",
	// php
	"php": "PHP",
	// windows command
	"cmd": "Windows Command Prompt",
	// linux shell
	"sh":   "Shell",
	"bash": "bash",
}

type SystemInfo struct {
	ARCH        string       `json:"arch"`
	OS          string       `json:"os"`
	Hostname    string       `json:"host_name"`
	NumCpu      int          `json:"num_cpu"`
	Executables []Executable `json:"executables"`
}

type Executable struct {
	Path        string `json:"path"`
	FileName    string `json:"file_name"`
	DisplayName string `json:"display_name"`
}

func GetLocalSystemInfo() (sysInfo SystemInfo, err error) {
	executables, err := GetExecutables()
	if err != nil {
		return sysInfo, err
	}
	hostname, err := os.Hostname()
	if err != nil {
		debug.PrintStack()
		return sysInfo, err
	}

	return SystemInfo{
		ARCH:        runtime.GOARCH,
		OS:          runtime.GOOS,
		NumCpu:      runtime.GOMAXPROCS(0),
		Hostname:    hostname,
		Executables: executables,
	}, nil
}

func GetSystemEnv(key string) string {
	return os.Getenv(key)
}

func GetPathValues() (paths []string) {
	pathEnv := GetSystemEnv("PATH")
	return strings.Split(pathEnv, ":")
}

func GetExecutables() (executables []Executable, err error) {
	pathValues := GetPathValues()

	cache := map[string]string{}

	for _, path := range pathValues {
		fileList, err := ioutil.ReadDir(path)
		if err != nil {
			log.Errorf(err.Error())
			debug.PrintStack()
			continue
		}

		for _, file := range fileList {
			displayName := executableNameMap[file.Name()]
			filePath := filepath.Join(path, file.Name())

			if cache[filePath] == "" {
				if displayName != "" {
					executables = append(executables, Executable{
						Path:        filePath,
						FileName:    file.Name(),
						DisplayName: displayName,
					})
				}
				cache[filePath] = filePath
			}
		}
	}
	return executables, nil
}
