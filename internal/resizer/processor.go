package resizer

import (
	"image"
	"io"

	"github.com/disintegration/imaging"
	"github.com/pkg/errors"
)

var _ ImageProcessor = (*processor)(nil)

type ImageProcessor interface {
	Decode(reader io.Reader) (image.Image, error)
	Encode(img image.Image, quality int, writer io.Writer) error
	Resize(img image.Image, width, height int) image.Image
}

type processor struct{}

func (i *processor) Decode(reader io.Reader) (image.Image, error) {
	img, err := imaging.Decode(reader)
	if err != nil {
		return nil, errors.Wrap(err, "image decode")
	}

	return img, nil
}

func (i *processor) Encode(img image.Image, quality int, writer io.Writer) error {
	if err := imaging.Encode(writer, img, imaging.JPEG, imaging.JPEGQuality(quality)); err != nil {
		return errors.Wrap(err, "image encode")
	}

	return nil
}

func (i *processor) Resize(img image.Image, width, height int) image.Image {
	return imaging.Resize(img, width, height, imaging.Lanczos)
}
