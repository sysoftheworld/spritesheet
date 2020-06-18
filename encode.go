package spritesheet

import (
	"errors"
	"image"
	"image/draw"
)

const (
	DefaultImgsPerRow = 5
)

var (
	ErrNoImages      = errors.New("no images passed to the encoder")
	ErrBadDimensions = errors.New("width and/or height of images passed is zero")
)

// NewAlpha returns new image.Alpha
func NewAlpha(r image.Rectangle) draw.Image {
	return image.NewAlpha(r)
}

// NewAlpha returns new image.Alpha16
func NewAlpha16(r image.Rectangle) draw.Image {
	return image.NewAlpha16(r)
}

// NewCMYK returns new image.CMYK
func NewCMYK(r image.Rectangle) draw.Image {
	return image.NewCMYK(r)
}

// NewGrey returns new image.Grey
func NewGray(r image.Rectangle) draw.Image {
	return image.NewGray(r)
}

// NewGray16 returns new image.Gray16
func NewGray16(r image.Rectangle) draw.Image {
	return image.NewGray16(r)
}

// NewNRGBA returns new image.NRGBA
func NewNRGBA(r image.Rectangle) draw.Image {
	return image.NewNRGBA(r)
}

// NewNRGBA64 returns new image.NRGBA64
func NewNRGBA64(r image.Rectangle) draw.Image {
	return image.NewNRGBA64(r)
}

// NewRGBA returns new image.NRGBA64
func NewRGBA(r image.Rectangle) draw.Image {
	return image.NewRGBA(r)
}

// NewRGBA64 returns new image.NewRGBA64
func NewRGBA64(r image.Rectangle) draw.Image {
	return image.NewRGBA64(r)
}

// EncodeOpts provides the encoder the parameters it needs to create a sprite sheet
type EncodeOpts struct {
	New        func(r image.Rectangle) draw.Image // what format you want the new image to be, defaults to RGBA
	ImgsPerRow int
}

// Encode takes a slice of images and based on the encode options will turn
// the images into a single sprite sheet.
// If the images are not all the same size, encode will take the max height
// and max width of any image and use that as its dimensions. For best look,
// all images should be the same size.
func Encode(images []image.Image, opts *EncodeOpts) (image.Image, error) {
	images = removeNilImages(images)
	if len(images) == 0 {
		return nil, ErrNoImages
	}
	if opts == nil {
		opts = &EncodeOpts{}
	}
	var (
		imgsPerRow    = imgsPerRow(opts.ImgsPerRow, len(images))
		mw, mh        = maxDimensions(images)
		width, height = sheetDimensions(imgsPerRow, len(images), mw, mh)
	)
	if width == 0 || height == 0 {
		return nil, ErrBadDimensions
	}

	var (
		row, column int
		sheet       = newSheet(width, height, opts.New)
	)
	for i := range images {
		var (
			x, y     = column * mw, row * mh                            // where we are at
			bounds   = images[i].Bounds()                               // current img bounds
			subImage = image.Rect(x, y, x+bounds.Max.X, y+bounds.Max.Y) // sub rectagle
		)
		draw.Draw(sheet, subImage, images[i], bounds.Min, draw.Over)

		column++
		if column > imgsPerRow-1 {
			row++
			column = 0
		}
	}
	return sheet, nil
}

// remove any nil interfaces
func removeNilImages(in []image.Image) (out []image.Image) {
	out = make([]image.Image, 0, len(in))
	for i := range in {
		if in[i] != nil {
			out = append(out, in[i])
		}
	}
	return
}

func imgsPerRow(imgsPerRow, numOfImgs int) int {
	if imgsPerRow == 0 {
		imgsPerRow = DefaultImgsPerRow
	}
	if numOfImgs < imgsPerRow {
		imgsPerRow = numOfImgs
	}
	return imgsPerRow
}

// sheetDimensions returns the dimensions of a new sprite sheet based on the number of images,
// images desired per row, and the max width and height of any of the images.
func sheetDimensions(imgsPerRow, numOfImgs, width, height int) (w int, h int) {
	if width == 0 || imgsPerRow == 0 {
		return // avoid div by 0
	}
	w = imgsPerRow * width
	rows := numOfImgs / imgsPerRow
	if numOfImgs%imgsPerRow != 0 {
		rows++
	}
	h = height * rows
	return
}

// newSheet creates a new sprite sheet ready to be "drawed" on based on desired widthXheight
// Defaults to using image.RGBA if no new function is provided.
func newSheet(width, height int, n func(rect image.Rectangle) draw.Image) draw.Image {
	want := image.Rect(0, 0, width, height)
	if n == nil { // default
		return image.NewRGBA(want)
	}
	return n(want)
}

// maxDimensions goes through a slice of images and returns the max width and height
func maxDimensions(images []image.Image) (w, h int) {
	if len(images) == 0 {
		return
	}
	for i := range images {
		if images[i] == nil {
			continue
		}
		var (
			size = images[i].Bounds().Size()
			x, y = size.X, size.Y
		)
		if x > w {
			w = x
		}
		if y > h {
			h = y
		}
	}
	return
}
