package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/khalieb/spritesheet"
)

var (
	in, out, op, dimensions string
	imgsPerRow              int
)

func init() {
	flag.StringVar(&in, "i", "", "input file(s); use dash seperated values")
	flag.StringVar(&out, "o", "", "output file")
	flag.StringVar(&op, "op", "encode", "decode or encode; default encode")
	flag.StringVar(&dimensions, "dim", "", "dimensions of image to decode")
	flag.IntVar(&imgsPerRow, "imgs-per-row", spritesheet.DefaultImgsPerRow, "number of images per row on encode")
}

func main() {
	flag.Parse()

	filenames, err := filenames(in)
	if err != nil {
		log.Fatal(err)
	}

	images, err := images(filenames)
	if err != nil {
		log.Fatal(err)
	}

	if len(images) == 0 {
		log.Println("no images to process")
		return
	}

	switch op {
	case "encode":
		err = encode(images, out, spritesheet.EncodeOpts{ImgsPerRow: imgsPerRow})
	case "decode":
		var w, h int
		w, h, err = parseDimensions(dimensions)
		if err != nil {
			break
		}
		err = decode(images[0], out, spritesheet.DecodeOpts{Width: w, Height: h})
	default:
		log.Fatalf("unknown operation %s", op)

	}

	if err != nil {
		log.Fatal(err)
	}
}

func filenames(in string) (out []string, err error) {
	// read from stdin if input is not provided
	var (
		rd    io.Reader
		delim byte = '\n'
	)
	if in != "" {
		rd = strings.NewReader(in)
		delim = ','
	} else {
		rd = os.Stdin
	}

	var (
		eof   bool
		bufRd = bufio.NewReader(rd)
	)
	for {
		line, err := bufRd.ReadString(delim)
		if err != nil && err != io.EOF {
			return nil, err
		}
		if err == io.EOF {
			eof = true
		}
		line = strings.TrimSpace(line)
		if len(line) == 0 {
			break
		}
		if line[len(line)-1] == delim {
			line = line[0 : len(line)-1]
		}

		out = append(out, line)
		if eof {
			break
		}
	}
	return
}

func images(in []string) ([]image.Image, error) {
	out := make([]image.Image, 0, len(in))
	for _, name := range in {
		file, err := os.Open(name)
		if err != nil {
			return nil, err
		}
		defer file.Close()

		format, err := format(name)
		if err != nil {
			return nil, err
		}

		image, err := decodeFile(file, format)
		if err != nil {
			return nil, err
		}
		out = append(out, image)
		file.Close()
	}
	return out, nil
}

func encode(images []image.Image, out string, opts spritesheet.EncodeOpts) error {
	sheet, err := spritesheet.Encode(images, &opts)
	if err != nil {
		return err
	}

	var w io.Writer
	if out == "" {
		w = os.Stdout
	} else {
		f, err := os.Create(out)
		if err != nil {
			return err
		}
		defer f.Close()
		w = f
	}
	format, err := format(out)
	if err != nil {
		return err
	}
	return encodeFile(w, sheet, format)
}

func decode(image image.Image, out string, opts spritesheet.DecodeOpts) error {
	images, err := spritesheet.Decode(image, &opts)
	if err != nil {
		return err
	}

	if out == "" {
		out = "sprite"
	}

	for i := range images {
		f, err := os.Create(fmt.Sprintf("%s-%d", out, i))
		if err != nil {
			return err
		}
		defer f.Close()

		format, err := format(out)
		if err != nil {
			return err
		}
		err = encodeFile(f, images[i], format)
		if err != nil {
			return err
		}
	}
	return nil
}

func format(filename string) (string, error) {
	ext := filepath.Ext(filename)
	if ext == "" {
		return "", fmt.Errorf("no extensions found in %s", filename)
	}
	return "jpeg", nil
}

func decodeFile(rd io.Reader, format string) (out image.Image, err error) {
	switch format {
	case "jpeg", "jpg":
		out, err = jpeg.Decode(rd)
	default:
		err = nil
	}
	return
}

func encodeFile(w io.Writer, image image.Image, format string) error {
	return jpeg.Encode(w, image, nil)
}

func parseDimensions(dim string) (width, height int, err error) {
	if dim == "" {
		err = errors.New("dimensions are required to decode")
	}
	if len(dim) > 20 {
		err = fmt.Errorf("these dimensions seem absurd (%s); we want nothing to do with them", dim)
	}
	if err != nil {
		return
	}

	var delim int
	for i := range dim {
		if dim[i] == 'x' || dim[i] == 'X' {
			delim = i
			break
		}
	}
	if delim == 0 {
		err = fmt.Errorf("could not parse dimensions: missing x or X in %s", dim)
	}
	if delim == len(dim) {
		err = fmt.Errorf("could not parse dimensions: x or X found at end of string in %s", dim)
	}
	if err != nil {
		return
	}
	w := dim[:delim]
	h := dim[delim+1 : len(dim)]

	width, err = strconv.Atoi(w)
	if err != nil {
		return
	}
	height, err = strconv.Atoi(h)
	if err != nil {
		return
	}
	return
}
