package spritesheet

import (
	"image"
	"image/draw"
)

// DecodeOpts provides the decoder with the necessary opts to split the spritesheet into seperate images
type DecodeOpts struct {
	New           func(r image.Rectangle) draw.Image // what format you want the new image to be, defaults to RGBA
	Width, Height int
}

// Decode takes in a image, assumed to be a spritesheet, and based on the options passed will
// chop up the spritesheet into seperate images.
// a width and height are needed to know the bounds of each image
func Decode(in image.Image, opts *DecodeOpts) ([]image.Image, error) {
	if in == nil {
		return nil, nil
	}
	if opts == nil || opts.Width == 0 || opts.Height == 0 {
		return nil, ErrBadDimensions
	}
	if opts.New == nil {
		opts.New = NewRGBA
	}

	var (
		row, column   int
		width, height = opts.Width, opts.Height
		out           []image.Image
		bounds        = in.Bounds()
	)
	for {
		var (
			min    = image.Point{X: (bounds.Min.X + column) * width, Y: (bounds.Min.Y + row) * height}
			max    = image.Point{X: min.X + width, Y: min.Y + height}
			subImg = image.Rectangle{Min: min, Max: max}
			newImg = opts.New(subImg)
		)

		column++
		if max.X >= bounds.Max.X {
			row++
			column = 0
		}

		if max.Y > bounds.Max.Y {
			break
		}
		draw.Draw(newImg, subImg, in, newImg.Bounds().Min, draw.Over)
		out = append(out, newImg)
	}
	return out, nil
}
