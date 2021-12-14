package main

import (
	"fmt"
	"os"
	"io/fs"
	"path/filepath"
	"net/http"
	"image"
	"image/jpeg"
	"image/png"
	"strings"
	"bytes"
	"log"
)

func checkErr(err error, path string, message string) {
	if err != nil {
		fmt.Printf("error: %s: %s\n", path, message)
		os.Exit(1)
	}
}

type convImg struct {
	in string
	out string
	oldPath string
	newPath string
	fileBytes []byte
}

func getImgInfo (path string) (string, []byte, error) {
	file, err := os.Open(path)
	checkErr(err, path, "could not open file")
	defer file.Close()
	fileBytes, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	contentType := http.DetectContentType(fileBytes)
	return contentType, fileBytes, file.Close()
}

func makeConvImg (in string, out string, oldPath string, fileBytes []byte) convImg {
	parts := strings.Split(oldPath, in)
	newPath := parts[0] + out
	return convImg{in, out, oldPath, newPath, fileBytes}
}

func (myImg convImg) encode() image.Image{
	var img image.Image
	var err error

	if myImg.in == "jpg" {
		img, err = jpeg.Decode(bytes.NewReader(myImg.fileBytes))
	} else if myImg.in == "png" {
		img, err = png.Decode(bytes.NewReader(myImg.fileBytes))
	}
	if (err != nil) {
		log.Fatal(err)
	}
	return img
}

func (myImg convImg) decode(img image.Image) bytes.Buffer {
	var buf bytes.Buffer
	var err error

	if myImg.out == "jpg" {
		err = jpeg.Encode(&buf, img, nil)
	} else if myImg.out == "png" {
		err = jpeg.Encode(&buf, img, nil)
	}
	if (err != nil) {
		log.Fatal(err)
	}
	return buf
}

func parseArgs() (string, string) {
	argc := len(os.Args)
	var in string
	var out string

	for i := 0; i < argc - 1; i++ {
		if strings.Contains(os.Args[i], "-i=") {
			in = strings.Trim(os.Args[i], "-i=")
		}else if strings.Contains(os.Args[i], "-o=") {
			out = strings.Trim(os.Args[i], "-o=")
		}
	}
	return in, out
}
func main() {
	var in string
	var out string

	if len(os.Args) != 4 {
		fmt.Println("error: invalid argument")
		os.Exit(1)
	}
	in, out = parseArgs()
	if in == out || (in != "png" && out != "png") || (in != "jpg" && out != "jpg") {
		fmt.Printf("Either %s or %s conversion is not available\n", in, out)
		os.Exit(1)
	}
	filepath.WalkDir(os.Args[3], func(path string, file fs.DirEntry, err error) error {
		checkErr(err, os.Args[3], "no such file or directory")
		if !file.IsDir() {
			contentType, fileBytes, err := getImgInfo(path)
			if err != nil {
				log.Fatal(err)
			}
			if contentType != "image/jpeg" && contentType != "image/png" {
				fmt.Printf("error: %s is not a valid file\n", path)
				return nil;
			}
			if in == "jpg" && contentType != "image/jpeg" || in == "png" && contentType != "image/png" {
				fmt.Printf("contentType: %s\n", contentType)
				fmt.Printf("error: -i option does not match %s image type\n", path)
				return nil;
			}
			myImg := makeConvImg(in, out, path, fileBytes)
			encodedImg := myImg.encode()
			buf := myImg.decode(encodedImg)
			file, err := os.Create(myImg.newPath)
			if (err != nil) {
				log.Fatal(err)
			}
			buf.WriteTo(file)
			file.Close()
		}
		return nil
	})
}

//buf := new(bytes.Buffer) if using new, you don't need to send the reference//
//pointers of type bytes.Buffer implement the io.Reader & io.Writer interfacies
