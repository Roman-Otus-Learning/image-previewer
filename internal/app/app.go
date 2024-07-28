package app

import (
	"context"
	"net/http"

	"github.com/Roman-Otus-Learning/image-previewer/internal/client"
	"github.com/Roman-Otus-Learning/image-previewer/internal/resizer"
	"github.com/pkg/errors"
)

var ErrRequestError = errors.New("request error")

type App interface {
	ResizeImage(ctx context.Context, url string, width, height int, headers http.Header) ([]byte, error)
}

func CreateResizerApp(client client.Client, resizer resizer.Resizer, imageQuality int) *ResizerApp {
	return &ResizerApp{client, resizer, imageQuality}
}

type ResizerApp struct {
	client       client.Client
	resizer      resizer.Resizer
	imageQuality int
}

func (a *ResizerApp) ResizeImage(
	ctx context.Context,
	url string,
	width, height int,
	headers http.Header,
) ([]byte, error) {
	rsp, err := a.client.GetWithHeaders(ctx, url, headers)
	if err != nil {
		return nil, errors.Wrap(err, "client get")
	}
	defer rsp.Body.Close()

	if rsp.StatusCode != http.StatusOK {
		return nil, ErrRequestError
	}

	result, err := a.resizer.Execute(rsp.Body, width, height, a.imageQuality)
	if err != nil {
		return nil, errors.Wrap(err, "resizer execute")
	}

	return result, err
}
