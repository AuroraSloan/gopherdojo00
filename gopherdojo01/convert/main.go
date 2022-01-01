package main

import (
	imgconv "convert/imgconv"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	in, out, dirPath := parseArgs()

	// Walk dirPath and run Convert on every item in dirPath that is not a directory
	err := filepath.WalkDir(dirPath, func(path string, file fs.DirEntry, err error) error {

		//check error from anonymous function
		checkErr(err, path, "no such file or directory")

		//convert image if the current file/path is not a directory
		if !file.IsDir() {
			_, err := imgconv.Convert(in, out, path)

			//Exit program and in case of error from Convert
			if err != nil {
				log.Fatal(err)
			}
		}
		//Return from anoynomous function
		return nil
	})

	//Exit program and in case of error from WalkDir
	if err != nil {
		log.Fatal(err)
	}
}

func checkErr(err error, path string, message string) {
	//Exit program if automatic funtion in WalkDir returned an error
	if err != nil {
		fmt.Printf("error: %s: %s\n", path, message)
		os.Exit(1)
	}
}

func parseArgs() (string, string, string) {
	argc := len(os.Args)
	var in string
	var out string

	//Exit if there are not 2 or 4 args
	//Return jpg for in and png for out in case of only one argument
	if argc != 2 && argc != 4 {
		exitConvert("error: invalid argument")
	} else if argc == 2 {
		return ".jpg", ".png", os.Args[1]
	}

	//decide in and out based on args
	in, out = setInOut()

	//Exit failure if in and out match, or if they neither jpg or png
	if in == out || (in != ".png" && out != ".png") || (in != ".jpg" && out != ".jpg") {
		fmt.Printf("Either %s or %s conversion is not available\n", in, out)
		os.Exit(1)
	}

	//return in type, out type, and directory path
	return in, out, os.Args[3]
}

func setInOut() (string, string) {
	var in string
	var out string

	//Exit if -i option arg is written not written correctly
	//If option is correct concatinate input type with '.'
	if !strings.Contains(os.Args[1], "-i=") {
		exitConvert("usage: ./convert -i=<type> -o=<type> directory")
	} else {
		in = "." + strings.Trim(os.Args[1], "-i=")
	}

	//Exit if -i option arg is written not written correctly.
	//If option is correct concatinate input type with '.'
	if !strings.Contains(os.Args[2], "-o=") {
		exitConvert("usage: ./convert -i=<type> -o=<type> directory")
	} else {
		out = "." + strings.Trim(os.Args[2], "-o=")
	}

	//return in type and out type
	return in, out
}

func exitConvert(message string) {
	//print message and exit program with failing status
	fmt.Println(message)
	os.Exit(1)
}
