package server

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/Roman-Otus-Learning/image-previewer/internal/app"
	"github.com/rs/zerolog/log"
)

const (
	pathPartsExpected  = 4
	pathPartsWidthIdx  = 1
	pathPartsHeightIdx = 2
	pathPartsURLIdx    = 3
	badRequestText     = "bad request"
)

var (
	ErrIncorrectRequestPath = errors.New("incorrect request path")
	ErrIncorrectWidth       = errors.New("incorrect width")
	ErrIncorrectHeight      = errors.New("incorrect height")
)

type Handler struct {
	app app.App
}

type request struct {
	width, height int
	url           string
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rq, err := parsePath(r.URL.Path)
	if err != nil {
		log.Error().Err(err).Msg("parse path " + r.URL.Path)
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	resized, err := h.app.ResizeImage(r.Context(), rq.url, rq.width, rq.height, r.Header)
	if err != nil {
		log.Error().Err(err).Msg("resize")
		w.WriteHeader(http.StatusBadGateway)
		_, _ = w.Write([]byte(badRequestText))
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(resized)
	log.Info().Msg("resized successfully: " + rq.url)
}

func parsePath(path string) (*request, error) {
	parts := strings.SplitN(path, `/`, pathPartsExpected)
	if len(parts) != pathPartsExpected {
		return nil, ErrIncorrectRequestPath
	}

	w, err := strconv.Atoi(parts[pathPartsWidthIdx])
	if err != nil {
		return nil, fmt.Errorf("%s: %w", parts[pathPartsWidthIdx], ErrIncorrectWidth)
	}

	h, err := strconv.Atoi(parts[pathPartsHeightIdx])
	if err != nil {
		return nil, fmt.Errorf("%s: %w", parts[pathPartsHeightIdx], ErrIncorrectHeight)
	}

	return &request{w, h, parts[pathPartsURLIdx]}, nil
}
