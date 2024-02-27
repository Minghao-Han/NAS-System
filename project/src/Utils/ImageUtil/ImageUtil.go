package ImageUtil

import (
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
)

type Decode func(reader io.Reader) (image.Image, error)

var Decoders = make(map[string]Decode)

func init() {
	Decoders[".jpg"] = jpeg.Decode
	Decoders[".jpeg"] = jpeg.Decode
	Decoders[".png"] = png.Decode
	Decoders[".gif"] = gif.Decode
}

func ImgDecode(reader io.Reader, imgType string) (image.Image, error) {
	decoder := Decoders[imgType]
	if decoder == nil {
		return nil, fmt.Errorf("unsupported image type")
	}
	return decoder(reader)
}
