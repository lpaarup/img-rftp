package main

import (
	"flag"
	"fmt"
	"image"
	"image/png"
	"io"
	"os"
	"path/filepath"

	"github.com/kevin-cantwell/dotmatrix"
	"github.com/lpaarup/img-rftp/pkg/client"
	"github.com/nfnt/resize"
	"github.com/pkg/errors"
)

var port = flag.Int("port", 69, "server's port")

func main() {
	flag.Parse()

	arguments := flag.Args()
	if len(arguments) == 0 {
		fmt.Fprintf(os.Stderr, "USAGE:\n\t tftprint [network images to print]\n")
		os.Exit(2)
	}

	c := client.New(fmt.Sprintf("127.0.0.1:%d", *port))
	for _, path := range arguments {

		r, err := c.Read(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "could not read the image from the server: %v", err)
			return
		}

		ext := filepath.Ext(path)
		if ext == "" {
			ext = ".png"
		}
		if err := print(r, ext[1:]); err != nil {
			fmt.Fprintf(os.Stderr, "could not print %s: %v", path, err)
		}
	}
}

func print(r io.Reader, format string) error {
	var img image.Image
	var err error
	switch {
	case format == "png":
		img, err = png.Decode(r)
		if err != nil {
			return errors.Wrap(err, "could not decode png image")
		}
	case format == "jpg" || format == "jpeg":
		img, _, err = image.Decode(r)
		if err != nil {
			return errors.Wrap(err, "could not decode jpg image")
		}
	}

	newImg := resize.Resize(160, 0, img, resize.Lanczos3)

	dotmatrix.Print(os.Stdout, newImg)
	fmt.Println()

	return nil
}
