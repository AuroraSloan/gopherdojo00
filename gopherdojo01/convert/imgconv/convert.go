/*
The MakeConversion function takes the in type, out type, and image path.
It will return a bool specifying whether the file is valid or not and an error in case of system error.

The only types available for conversion currently are jpg and png.


*/
package imgconv

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"strings"
)

func Convert(in string, out string, path string) (bool, error) {

	// Reads image file and returns content type and image file bytes
	contentType, fileBytes, err := getImgInfo(path)
	if err != nil {
		return false, err
	}

	// Compares in type and out type with content type. Returns boolean specifying whether
	// or not the image is valid or not
	if validConv := checkContentType(in, out, contentType, path); !validConv {
		return false, nil
	}

	// Create convInfo type containing relevent path and image information
	myImg := gatherInfo(in, out, path, fileBytes)

	// Encodes the specified image to to an image.Image
	decodedImg := myImg.decode()

	// Decodes the specified type and writes it to a buffer
	buf := myImg.encode(decodedImg)

	// Writes buffer from decoded image to the specified new image file
	// and returns it
	return createNewImage(buf, myImg.NewPath)
}

func getImgInfo(path string) (string, []byte, error) {
	// Open file for reading
	file, err := os.Open(path)
	if err != nil {
		return "", nil, err
	}

	// Read file to fileBytes
	fileBytes, err := os.ReadFile(path)
	if err != nil {
		file.Close()
		return "", nil, err
	}

	// Detect content type and set to var contentType
	contentType := http.DetectContentType(fileBytes)

	// Returns contentType, fileBytes and close() status of file
	return contentType, fileBytes, file.Close()
}

func checkContentType(in string, out string, contentType string, path string) bool {

	// Returns false if something other than jpeg or png is returned
	if contentType != "image/jpeg" && contentType != "image/png" {
		fmt.Printf("error: %s is not a valid file\n", path)
		return false
	}

	// Returns false if file does not end as .jpg or .png
	if parts := strings.Split(path, in); len(parts) != 2 || parts[1] == out {
		return false
	}
	return true
}

func createNewImage(buf bytes.Buffer, newPath string) (bool, error) {

	// Open file for writing
	file, err := os.Create(newPath)
	if err != nil {
		return false, err
	}

	// Write image buffer to file and close file
	_, err = buf.WriteTo(file)
	defer file.Close()
	if err != nil {
		return false, err
	}

	return true, file.Close()
}
