/*
Package imgconv converts images from one type to another.
*/
package imgconv

import (
	"bytes"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"log"
	"strings"
)

// Type containing relevent path and image information
type convInfo struct {
	in        string
	out       string
	oldPath   string
	NewPath   string
	fileBytes []byte
}

// Creates convInfo type
func gatherInfo(in string, out string, oldPath string, fileBytes []byte) convInfo {

	// Create new file name
	parts := strings.Split(oldPath, in)
	newPath := parts[0] + out

	// reutrn convInfo struct
	return convInfo{in, out, oldPath, newPath, fileBytes}
}

// Encodes the specified image to an image.Image
func (myImg convInfo) decode() image.Image {
	var img image.Image
	var err error

	// Decode in image file bytes to img
	switch myImg.in {
	case ".jpg":
		img, err = jpeg.Decode(bytes.NewReader(myImg.fileBytes))
	case ".png":
		img, err = png.Decode(bytes.NewReader(myImg.fileBytes))
	}
	if err != nil {
		log.Fatal(err)
	}

	// Draw white background underneath the png image to make up for lack
	// of transparency in jpegs
	if myImg.in == ".png" {
		newImg := image.NewRGBA(img.Bounds())
		draw.Draw(newImg, newImg.Bounds(), &image.Uniform{color.White}, image.Point{}, draw.Src)
		draw.Draw(newImg, newImg.Bounds(), img, img.Bounds().Min, draw.Over)
		return newImg
	}

	return img
}

// Decodes the specified type and writes it to a buffer
func (myImg convInfo) encode(img image.Image) bytes.Buffer {
	var buf bytes.Buffer
	var err error

	// Encode previously decoded img to buf and return
	switch myImg.out {
	case ".jpg":
		err = jpeg.Encode(&buf, img, nil)
	case ".png":
		err = png.Encode(&buf, img)
	}
	if err != nil {
		log.Fatal(err)
	}
	return buf
}