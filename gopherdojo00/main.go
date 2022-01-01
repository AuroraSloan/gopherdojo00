package main

import (
	"fmt"
	"os"
	"io/fs"
	"path/filepath"
	"log"
	"strings"
	imgconv "github.com/AuroraSloan/gopherdojo00/Imgconv"
)

func main() {
	in, out, dirPath := parseArgs()

	// Walk dirPath and run MakeConversion on every item in dirPath that is not a directory
	err := filepath.WalkDir(dirPath, func(path string, file fs.DirEntry, err error) error {
		checkErr(err, path, "no such file or directory")
		if !file.IsDir() {
			err := imgconv.MakeConversion(in, out, path)
			if err != nil {
				log.Fatal(err)
			}
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
}

func checkErr(err error, path string, message string) {
	if err != nil {
		fmt.Printf("error: %s: %s\n", path, message)
		os.Exit(1)
	}
}

func parseArgs() (string, string, string) {
	argc := len(os.Args)
	var in string
	var out string

	//check number of arguments
	if argc != 2 && argc != 4 {
		fmt.Println("error: invalid argument")
		os.Exit(1)
	}else if argc == 2 {
		return ".jpg", ".png", os.Args[1]
	}

	//check arguments -i and -o options match usage
	if !strings.Contains(os.Args[1], "-i=") {
		fmt.Println("usage: ./convert -i=<type> -o=<type> directory")
		os.Exit(1)
	} else {
		in = "." + strings.Trim(os.Args[1], "-i=")
	}
	if !strings.Contains(os.Args[2], "-o=") {
		fmt.Println("usage: ./convert -i=<type> -o=<type> directory")
		os.Exit(1)
	} else {
		out = "." + strings.Trim(os.Args[2], "-o=")
	}

	//confirm in and out are not the same and are either jpg or png
	if in == out || (in != ".png" && out != ".png") || (in != ".jpg" && out != ".jpg") {
		fmt.Printf("Either %s or %s conversion is not available\n", in, out)
		os.Exit(1)
	}
	return in, out, os.Args[3]
}
