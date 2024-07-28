package resizer

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestResizer(t *testing.T) {
	resizer := CreateResizer()
	t.Run("success", func(t *testing.T) {
		imgPath := "../../sample/picture.jpg"
		expectedWidth := 100
		expectedHeight := 100
		imgData, err := os.ReadFile(imgPath)
		assert.NoError(t, err)

		reader := bytes.NewReader(imgData)
		resizedImage, err := resizer.Execute(reader, expectedWidth, expectedHeight, 70)
		assert.NoError(t, err)

		processor := &processor{}
		img, err := processor.Decode(bytes.NewReader(resizedImage))
		assert.NoError(t, err)

		resultWidth, resultHeight := img.Bounds().Dx(), img.Bounds().Dy()
		require.Equal(t, resultWidth, expectedWidth)
		require.Equal(t, resultHeight, expectedHeight)
	})

	t.Run("unknown file format", func(t *testing.T) {
		_, err := resizer.Execute(new(bytes.Buffer), 100, 100, 70)
		assert.Error(t, err, "failed to scale testcase img file")
	})
}
