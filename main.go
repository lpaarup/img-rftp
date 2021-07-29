package main

import (
	"fmt"
	"image/png"
	"io"
	"os"

	"github.com/kevin-cantwell/dotmatrix"
	"github.com/lpaarup/img-rftp/client"
	"github.com/nfnt/resize"
	"github.com/pkg/errors"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "missign paths of images to cat")
		os.Exit(2)
	}
	for _, path := range os.Args[1:] {
		c := client.New("127.0.0.1:8080")

		r, err := c.Read(os.Args[1])

		if err != nil {
			fmt.Fprintf(os.Stderr, "could not read the image from the server: %v", err)
		}

		if err := print(r); err != nil {
			fmt.Fprintf(os.Stderr, "could not print %s: %v", path, err)
		}
	}
}

func print(r io.Reader) error {
	img, err := png.Decode(r)
	if err != nil {
		return errors.Wrap(err, "could not decode image")
	}

	newImg := resize.Resize(160, 0, img, resize.Lanczos3)

	dotmatrix.Print(os.Stdout, newImg)

	return nil
}
