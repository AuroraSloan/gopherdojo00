/*
The MakeConversion function takes the in type, out type, and image path.
It will return a bool specifying whether the file is valid or not and an error in case of system error.

The only types available for conversion currently are jpg and png.


*/
package imgconv

import (
	"fmt"
	"os"
	"log"
	"net/http"
	"strings"
	"bytes"
)

func MakeConversion(in string, out string, path string) error {
	contentType, fileBytes, err := getImgInfo(path)
	if err != nil {
		return err
	}
	if validConv := checkContentType(in, out, contentType, path); !validConv {
		return nil
	}

	// Create convInfo type containing relevent path and image information
	myImg := gatherInfo(in, out, path, fileBytes)

	// Encodes the specified image to to an image.Image
	decodedImg := myImg.decode()

	// Decodes the specified type and writes it to a buffer
	buf := myImg.encode(decodedImg)

	// Writes buffer from decoded image to the specified new image file
	return createNewImage(buf, myImg.NewPath)
}

// Reads image and returns content type and image file bytes
func getImgInfo (path string) (string, []byte, error) {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	fileBytes, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	contentType := http.DetectContentType(fileBytes)
	return contentType, fileBytes, file.Close()
}

// Compares in type and out type with content type. Returns boolean specifying whether
// or not the image is valid or not
func checkContentType(in string, out string, contentType string, path string) bool {
	if contentType != "image/jpeg" && contentType != "image/png" {
		fmt.Printf("error: %s is not a valid file\n", path)
		return false;
	}
	if in == "jpg" && contentType != "image/jpeg" || in == "png" && contentType != "image/png" {
		fmt.Printf("error: -i option does not match %s image type\n", path)
		return false;
	}
	if parts := strings.Split(path, in); len(parts) != 2 || parts[1] == out {
		return false;
	}
	return true
}

// Writes buffer from decoded image to the specified new image file.
func createNewImage(buf bytes.Buffer, newPath string) error {
	file, err := os.Create(newPath)
	if (err != nil) {
		return err
	}
	_, err = buf.WriteTo(file)
	file.Close()
	if err != nil {
		return err
	}
	return nil
}
