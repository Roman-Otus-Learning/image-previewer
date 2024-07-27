package app

import (
	"context"
	"github.com/pkg/errors"
	"net/http"

	"github.com/Roman-Otus-Learning/image-previewer/internal/client"
	"github.com/Roman-Otus-Learning/image-previewer/internal/resizer"
)

var ErrRequestError = errors.New("request error")

type App interface {
	Resize(ctx context.Context, url string, width, height int, headers http.Header) ([]byte, error)
}

func CreateResizerApp(client client.Client, resizer resizer.Resizer) *ResizerApp {
	return &ResizerApp{client, resizer}
}

type ResizerApp struct {
	client  client.Client
	resizer resizer.Resizer
}

func (a *ResizerApp) Resize(ctx context.Context, url string, width, height int, headers http.Header) ([]byte, error) {
	rsp, err := a.client.GetWithHeaders(ctx, url, headers)
	if err != nil {
		return nil, errors.Wrap(err, "resizer get")
	}
	defer rsp.Body.Close()

	if rsp.StatusCode != http.StatusOK {
		return nil, ErrRequestError
	}

	result, err := a.resizer.Resize(rsp.Body, width, height)
	if err != nil {
		return nil, errors.Wrap(err, "resizer resize")
	}

	return result, err
}
