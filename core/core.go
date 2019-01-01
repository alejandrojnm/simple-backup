package core

import (
	"fmt"
	"io/ioutil"
	_ "log"
	"os"
	"path/filepath"
	"time"
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

func RemoveOldFile(dir string, days float64) {
	/*
		Function to remove all file in dir older than *days
	*/

	//Read dir to search old files
	files, _ := ioutil.ReadDir(dir)

	for _, f := range files {
		fi, err := os.Stat(dir + f.Name())
		if err != nil {
			fmt.Println(err)
		}

		// Calculate the difference between now and ModTime
		now := time.Now()
		currTime := fi.ModTime()
		diff := now.Sub(currTime)

		//if the file ModTime is larger than days*24 then delete file
		if days*24 < diff.Hours() {
			err := os.RemoveAll(dir + f.Name())
			if err != nil {
				fmt.Println(err)
			}

		}
	}
}
