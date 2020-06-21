package spritesheet

import (
	"image"
	"image/color"
	"math/rand"
	"testing"
)

func TestDecode(t *testing.T) {
	tests := []struct {
		name          string
		width, height int // for the image we are creating
		points        []image.Point
		opts          *DecodeOpts
		err           error
	}{
		{name: "nil opts", width: 10, height: 10, points: nil, opts: nil, err: ErrBadDimensions},
		{name: "bad dimensions", width: 10, height: 10, points: nil, opts: &DecodeOpts{}, err: ErrBadDimensions},
		{
			name:  "50x20",
			width: 100, height: 20,
			points: []image.Point{
				{X: 5, Y: 5}, {X: 15, Y: 5}, {X: 25, Y: 5}, {X: 35, Y: 5}, {X: 45, Y: 5}, // row 1
				{X: 5, Y: 15}, {X: 15, Y: 15}, {X: 25, Y: 15}, {X: 35, Y: 15}, {X: 45, Y: 15}, // row 2
			},
			opts: &DecodeOpts{Width: 10, Height: 10},
			err:  nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var (
				r, g, b, a = uint8(rand.Intn(255)), uint8(rand.Intn(255)), uint8(rand.Intn(255)), uint8(rand.Intn(255))
				color      = color.RGBA{R: r, G: g, B: b, A: a}
				img        = buildTestImage(test.width, test.height, color, test.points...)
			)
			out, err := Decode(img, test.opts)
			if err != test.err {
				t.Errorf("got %v; wanted %v", err, test.err)
				return
			}
			if err != nil {
				return
			}
			for i, img := range out {
				var (
					bounds = img.Bounds()
				)
				// n is small just loop
				for _, point := range test.points {
					if !inside(point, bounds) {
						continue
					}
					c := img.At(point.X, point.Y)
					if c != color {
						t.Errorf("got %v; wanted %v on image %d", c, color, i+1)
					}
				}
			}
		})
	}
}

// inside: check to see if point is inside a rectangle
func inside(point image.Point, rec image.Rectangle) bool {
	var min, max = rec.Min, rec.Max
	return min.X <= point.X && min.Y <= point.Y && max.X >= point.X && max.Y >= point.Y
}

func buildTestImage(width, height int, color color.Color, points ...image.Point) image.Image {
	img := NewRGBA(image.Rect(0, 0, width, height))
	for _, point := range points {
		img.Set(point.X, point.Y, color) // assume its within the bounds
	}
	return img
}
