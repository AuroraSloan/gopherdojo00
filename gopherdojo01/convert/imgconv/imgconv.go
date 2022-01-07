/*
Package imgconv converts images from one type to another.
*/
package imgconv

import (
	"bytes"
	"errors"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"strings"
)

// Type containing relevant path and image information
type convInfo struct {
	in        string
	out       string
	oldPath   string
	NewPath   string
	fileBytes []byte
}

// Creates convInfo type
func gatherInfo(in string, out string, oldPath string, fileBytes []byte, contentType string) (convInfo, error) {

	var imgTypes = map[string]string{
		"image/jpeg": ".jpg",
		"image/png":  ".png",
		"image/bmp":  ".bmp",
		"image/tiff": ".tiff",
	}
	if imgTypes[contentType] != in {
		return convInfo{"", "", "", "", nil}, errors.New("-i= type does not match image content type")
	}
	// Create new file name
	parts := strings.Split(oldPath, in)
	newPath := parts[0] + out

	// return convInfo struct
	return convInfo{in, out, oldPath, newPath, fileBytes}, nil
}

// Encodes the specified image to an image.Image
func (myImg convInfo) decode() (image.Image, error) {
	var img image.Image
	var err error

	// Decode in image file bytes to img
	switch myImg.in {
	case ".jpg":
		img, err = jpeg.Decode(bytes.NewReader(myImg.fileBytes))
	case ".png":
		img, err = png.Decode(bytes.NewReader(myImg.fileBytes))
	default:
		err = errors.New("-i= type is not supported")
	}
	if err != nil {
		return nil, err
	}

	// Draw white background underneath the png image to make up for lack
	// of transparency in jpegs
	if myImg.in == ".png" {
		newImg := image.NewRGBA(img.Bounds())
		draw.Draw(newImg, newImg.Bounds(), &image.Uniform{color.White}, image.Point{}, draw.Src)
		draw.Draw(newImg, newImg.Bounds(), img, img.Bounds().Min, draw.Over)
		return newImg, nil
	}

	return img, nil
}

// Decodes the specified type and writes it to a buffer
func (myImg convInfo) encode(img image.Image) (bytes.Buffer, error) {
	var buf bytes.Buffer
	var err error

	// Encode previously decoded img to buf and return
	switch myImg.out {
	case ".jpg":
		err = jpeg.Encode(&buf, img, nil)
	case ".png":
		err = png.Encode(&buf, img)
	default:
		err = errors.New("-o= type is not supported")
	}
	return buf, err
}
