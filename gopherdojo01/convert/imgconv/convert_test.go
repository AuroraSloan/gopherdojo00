package imgconv_test

import (
	imgconv "convert/imgconv"
	"errors"
	"net/http"
	"os"
	"strings"
	"testing"
)

type imagePath string

func (name imagePath) changePath(oldT string, newT string) string {
	parts := strings.Split(string(name), oldT)
	newPath := parts[0] + newT
	return newPath
}

func TestConvert(t *testing.T) {
	errA := errors.New("-i= type does not match image content type")
	errB := errors.New("-i= type is not supported")
	errC := errors.New("-o= type is not supported")
	moderr := os.Chmod("../testdata/noperms.jpg", 0000)
	if moderr != nil {
		t.Fatal(moderr)
	}
	moderr = os.Chmod("../testdata/cantopen.png", 0000)
	if moderr != nil {
		t.Fatal(moderr)
	}
	t.Parallel()
	var imgConversions = []struct {
		in, out, path string
		want          error
	}{
		{".jpg", ".png", "../testdata/images/jpgs/gopher.jpg", nil},
		{".png", ".jpg", "../testdata/images/pngs/tokyotower.png", nil},
		{".png", ".jpg", "../testdata/images/jpgs/gopher.jpg", errA},
		{".jpg", ".png", "../testdata/images/pngs/tokyotower.png", errA},
		{".jpg", ".png", "convert_test.go", errA},
		{".bmp", ".jpg", "../testdata/sample.bmp", errB},
		{".png", ".tiff", "../testdata/images/pngs/tokyotower.png", errC},
	}
	for _, ic := range imgConversions {
		t.Run("Convert", func(t *testing.T) {
			err := imgconv.Convert(ic.in, ic.out, ic.path)
			if err != nil && err.Error() != ic.want.Error() {
				t.Error(err)
			}
			if err == nil {
				oldPath := imagePath(ic.path)
				newPath := oldPath.changePath(ic.in, ic.out)
				contentType := checkOutputFile(t, newPath)
				if contentType != ic.out {
					t.Error("Output image does not match requested output")
				}
			}
		})
	}
	t.Run("no perms", func(t *testing.T) {
		err := imgconv.Convert(".jpg", ".png", "../testdata/noperms.jpg")
		if err == nil {
			t.Error("Should return err.")
		}
	})
	t.Run("Invalid file", func(t *testing.T) {
		err := imgconv.Convert(".jpg", ".png", "../testdata/doesnotexist")
		if err == nil {
			t.Error("Should return err.")
		}
	})
	t.Run("can't create", func(t *testing.T) {
		err := imgconv.Convert(".jpg", ".png", "../testdata/cantopen.jpg")
		if err == nil {
			t.Error("Should return err.")
		}
	})
}

func checkOutputFile(t *testing.T, path string) string {
	t.Helper()

	file, err := os.Open(path)
	if err != nil {
		t.Fatal("sys error opening file")
	}
	fileBytes, err := os.ReadFile(path)
	if err != nil {
		file.Close()
		t.Fatal("sys error reading file")
	}
	contentType := http.DetectContentType(fileBytes)
	defer file.Close()
	if contentType == "image/jpeg" {
		return ".jpg"
	} else {
		return ".png"
	}
}
