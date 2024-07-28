package resizer

import (
	"bytes"
	"io"

	"github.com/pkg/errors"
)

var _ Resizer = (*ImageResizer)(nil)

type Resizer interface {
	Execute(i io.Reader, width, height, quality int) ([]byte, error)
}

func CreateResizer() *ImageResizer {
	return &ImageResizer{
		&processor{},
	}
}

type ImageResizer struct {
	processor ImageProcessor
}

func (r *ImageResizer) Execute(reader io.Reader, width, height, quality int) ([]byte, error) {
	img, err := r.processor.Decode(reader)
	if err != nil {
		return nil, errors.Wrap(err, "resizer decode")
	}

	img = r.processor.Resize(img, width, height)
	buffer := new(bytes.Buffer)
	if err := r.processor.Encode(img, quality, buffer); err != nil {
		return nil, errors.Wrap(err, "resizer encode")
	}

	return buffer.Bytes(), nil
}
