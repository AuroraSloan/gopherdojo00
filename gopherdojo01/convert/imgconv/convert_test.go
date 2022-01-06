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
		want          bool
	}{
		{".jpg", ".png", "../testdata/images/jpgs/gopher.jpg", true},
		{".png", ".jpg", "../testdata/images/pngs/tokyotower.png", true},
		{".png", ".jpg", "../testdata/images/jpgs/gopher.jpg", false},
		{".jpg", ".png", "../testdata/images/pngs/tokyotower.png", false},
		{".jpg", ".png", "convert_test.go", false},
	}
	for _, ic := range imgConversions {
		t.Run("Convert", func(t *testing.T) {
			got, err := imgconv.Convert(ic.in, ic.out, ic.path)
			if err != nil {
				t.Fatal(err)
			} else if got != ic.want {
				t.Error("Convert output is incorrect")
			}
			if got {
				oldPath := imagePath(ic.path)
				newPath := oldPath.changePath(ic.in, ic.out)
				contentType := checkOutputFile(t, newPath)
				if contentType != ic.out {
					t.Error("Output image does not match requested output")
				}
			}
		})
	}
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
