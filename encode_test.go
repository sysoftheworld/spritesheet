package spritesheet

import (
	"image"
	"testing"
)

func TestSheetDimensions(t *testing.T) {
	tests := []struct {
		name                  string
		numOfImgs, imgsPerRow int
		inWidth, inHeight     int
		outWidth, outHeight   int
	}{
		{name: "empty"},
		{name: "zero height", imgsPerRow: 10, numOfImgs: 5, inWidth: 10, inHeight: 0, outWidth: 50, outHeight: 0},
		{name: "no remainder", imgsPerRow: 10, numOfImgs: 10, inWidth: 10, inHeight: 10, outWidth: 100, outHeight: 10},
		{name: "add a row", imgsPerRow: 10, numOfImgs: 11, inWidth: 10, inHeight: 10, outWidth: 100, outHeight: 20},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.imgsPerRow = imgsPerRow(test.imgsPerRow, test.numOfImgs) // set defaults
			w, h := sheetDimensions(test.imgsPerRow, test.numOfImgs, test.inWidth, test.inHeight)
			if w != test.outWidth {
				t.Errorf("want width %d, got %d", test.outWidth, w)
			}
			if h != test.outHeight {
				t.Errorf("want height %d, got %d", test.outHeight, h)
			}
		})
	}
}

func TestMaxDimensions(t *testing.T) {
	tests := []struct {
		name          string
		in            []image.Image
		width, height int
	}{
		{name: "nil images", in: nil, width: 0, height: 0},
		{name: "nil image", in: []image.Image{nil, nil}, width: 0, height: 0},
		{
			name: "same size",
			in: []image.Image{
				image.NewAlpha(image.Rect(0, 0, 100, 100)),
				image.NewAlpha(image.Rect(0, 0, 100, 100)),
			},
			width: 100, height: 100,
		},
		{
			name: "different sizes",
			in: []image.Image{
				image.NewAlpha(image.Rect(0, 0, 100, 100)),
				image.NewAlpha(image.Rect(0, 0, 101, 101)),
			},
			width: 101, height: 101,
		},
		{
			name: "one from each",
			in: []image.Image{
				image.NewAlpha(image.Rect(0, 0, 200, 100)),
				image.NewAlpha(image.Rect(0, 0, 101, 150)),
			},
			width: 200, height: 150,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			w, h := maxDimensions(test.in)
			if w != test.width {
				t.Errorf("want width %d, got %d", test.width, w)
			}
			if h != test.height {
				t.Errorf("want height %d, got %d", test.height, h)
			}
		})
	}
}
