package services

import (
	"bytes"
	"errors"
	"github.com/nfnt/resize"
	"image"
	"image/jpeg"
	"image/png"
)

type NFNTResizer struct {
}

func (r *NFNTResizer) Resize(body []byte, imgType string, scales []float64) ([]image.Image, error) {
	var img image.Image
	var err error
	switch imgType {
	case "image/png":
		img, err = png.Decode(bytes.NewReader(body))
		if err != nil {
			return nil, err
		}
	case "image/jpeg":
		img, err = jpeg.Decode(bytes.NewReader(body))
		if err != nil {
			return nil, err
		}
	default:
		return nil, errors.New("can't decode this type of image")
	}

	w := float64(img.Bounds().Dx())
	h := float64(img.Bounds().Dy())
	imgs := make([]image.Image, len(scales))

	for i, s := range scales {
		imgs[i] = resize.Resize(uint(w*s), uint(h*s), img, resize.Lanczos3)
	}
	return imgs, nil

}
