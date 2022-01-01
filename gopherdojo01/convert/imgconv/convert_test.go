package imgconv_test

import (
	imgconv "convert/imgconv"
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
	t.Parallel()
	var imgConversions = []struct {
		in, out, path string
	}{
		{".jpg", ".png", "../testdata/images/jpgs/gopher.jpg"},
		{".png", ".jpg", "../testdata/images/pngs/tokyotower.png"},
	}
	for _, ic := range imgConversions {
		t.Run("Valid Convert", func(t *testing.T) {
			isValid, err := imgconv.Convert(ic.in, ic.out, ic.path)
			if err != nil {
				t.Fatal(err)
			} else if !isValid {
				t.Error("Should return isValid.")
			}
			oldPath := imagePath(ic.path)
			newPath := oldPath.changePath(ic.in, ic.out)
			contentType := checkOutputFile(t, newPath)
			if contentType != ic.out {
				t.Error("Output image does not match requested output")
			}
		})
		t.Run("Invalid in / out", func(t *testing.T) {
			isValid, err := imgconv.Convert(ic.out, ic.in, ic.path)
			if err != nil {
				t.Fatal(err)
			} else if isValid {
				t.Error("Should return !isValid.")
			}
		})
	}
	t.Run("Invalid image type", func(t *testing.T) {
		isValid, err := imgconv.Convert(".jpg", ".png", "convert_test.go")
		if err != nil {
			t.Fatal(err)
		} else if isValid {
			t.Error("Should return !isValid.")
		}
	})
	t.Run("Invalid file", func(t *testing.T) {
		_, err := imgconv.Convert(".jpg", ".png", "../testdata/doesnotexist")
		if err == nil {
			t.Fatal("Should return err.")
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
