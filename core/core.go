package core

import (
	"fmt"
	_ "log"
	"path/filepath"
)

func FindFile(targetDir string, pattern []string) ([]string, error) {
	/*
		Function to find file base on regx
	*/
	matches, err := filepath.Glob(targetDir + pattern[0])

	if err != nil {
		fmt.Println(err)
	}

	return matches, err
}
