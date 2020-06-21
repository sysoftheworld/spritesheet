# Spritesheet
Spritesheet is an Encoder and Decoder to build spritesheets from a list of images or seperate a spritesheet into multiple images. 

## Encoder
The encoder takes in a slice of images and some options. If no options are provided, defaults will be used. 

```Go
// EncodeOpts provides the encoder the parameters it needs to create a spritesheet
type EncodeOpts struct {
	New        func(r image.Rectangle) draw.Image // what format you want the new image to be, defaults to RGBA
	ImgsPerRow int // Default to 5
}

// Encode takes a slice of images and based on the encode options will turn
// the images into a single sprite sheet.
// If the images are not all the same size, encode will take the max height
// and max width of any image and use that as its dimensions. For best look,
// all images should be the same size.
func Encode(images []image.Image, opts *EncodeOpts) (image.Image, error) {
```

## Decoder 
The decoder takes a spritesheet and based on width and height chops the spritesheet up into a slice of seperate images. 

```Go
// DecodeOpts provides the decoder with the necessary opts to split the spritesheet into seperate images
type DecodeOpts struct {
	New           func(r image.Rectangle) draw.Image // what format you want the new image to be, defaults to RGBA
	Width, Height int
}

// Decode takes in a image, assumed to be a spritesheet, and based on the options passed will
// chop up the spritesheet into seperate images.
// a width and height are needed to know the bounds of each image
func Decode(in image.Image, opts *DecodeOpts) ([]image.Image, error) {
```

## TODO
- Finish the cli
