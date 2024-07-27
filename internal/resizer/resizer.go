package resizer

import (
	"bytes"
	"github.com/pkg/errors"
	"io"
)

var _ Resizer = (*ImageResizer)(nil)

type Resizer interface {
	Resize(i io.Reader, width, height int) ([]byte, error)
}

func CreateResizer() *ImageResizer {
	return &ImageResizer{
		&processor{},
	}
}

type ImageResizer struct {
	processor ImageProcessor
}

func (r *ImageResizer) WithProcessor(processor ImageProcessor) *ImageResizer {
	r.processor = processor

	return r
}

func (r *ImageResizer) Resize(reader io.Reader, width, height int) ([]byte, error) {
	img, err := r.processor.Decode(reader)
	if err != nil {
		return nil, errors.Wrap(err, "resizer decode")
	}

	img = r.processor.Resize(img, width, height)
	buffer := new(bytes.Buffer)
	if err := r.processor.Encode(img, buffer); err != nil {
		return nil, errors.Wrap(err, "resizer encode")
	}

	return buffer.Bytes(), nil
}
