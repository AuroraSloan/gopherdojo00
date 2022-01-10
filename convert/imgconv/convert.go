/*
The MakeConversion function takes the in type, out type, and image path.
It will return a bool specifying whether the file is valid or not and an error in case of system error.

The only types available for conversion currently are jpg and png.


*/
package imgconv

import (
	"bytes"
	"net/http"
	"os"
)

func Convert(in string, out string, path string) error {

	// Reads image file and returns content type and image file bytes
	contentType, fileBytes, err := getImgInfo(path)
	if err != nil {
		return err
	}

	// Create convInfo type containing relevant path and image information
	myImg, err := gatherInfo(in, out, path, fileBytes, contentType)
	if err != nil {
		return err
	}

	// Encodes the specified image to to an image.Image
	decodedImg, err := myImg.decode()
	if err != nil {
		return err
	}

	// Decodes the specified type and writes it to a buffer
	buf, err := myImg.encode(decodedImg)
	if err != nil {
		return err
	}
	// Writes buffer from decoded image to the specified new image file
	// and returns it
	return createNewImage(buf, myImg.NewPath)
}

func getImgInfo(path string) (string, []byte, error) {

	// Read file to fileBytes
	fileBytes, err := os.ReadFile(path)
	if err != nil {
		return "", nil, err
	}

	// Detect content type and set to var contentType
	contentType := http.DetectContentType(fileBytes)

	// Returns contentType, fileBytes and close() status of file
	return contentType, fileBytes, nil
}

func createNewImage(buf bytes.Buffer, newPath string) error {

	// Open file for writing
	file, err := os.Create(newPath)
	if err != nil {
		return err
	}

	// Write image buffer to file and close file
	_, err = buf.WriteTo(file)
	if err != nil {
		file.Close()
		return err
	}

	return file.Close()
}
