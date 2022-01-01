package main

import (
	"io"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		os.Exit(1)
	}
	for i, filename := range os.Args {
		if i == 0 {
			continue
		}
		if filename == "-" {
			print_file(os.Stdin, os.Stdout)
		} else {
			src, err := os.Open(filename)
			if err != nil {
				os.Exit(1)
			}
			defer src.Close()
			print_file(src, os.Stdout)
		}
	}
}

func print_file(src io.Reader, dst io.Writer) {
	_, err := io.Copy(dst, src)
	if err != nil {
		os.Exit(1)
	}
}
