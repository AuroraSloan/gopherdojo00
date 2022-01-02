package main

import (
	"io"
	"os"
)

func main() {
	exit_status := 0
	if len(os.Args) < 2 {
		print_file(os.Stdin, os.Stdout)
	}
	for _, filename := range os.Args[1:] {
		if filename == "-" {
			print_file(os.Stdin, os.Stdout)
		} else {
			src, err := os.Open(filename)
			if err != nil {
				exit_status = 1
				continue
			}
			defer src.Close()
			print_file(src, os.Stdout)
		}
	}
	os.Exit(exit_status)
}

func print_file(src io.Reader, dst io.Writer) {
	_, err := io.Copy(dst, src)
	if err != nil {
		os.Exit(1)
	}
}
