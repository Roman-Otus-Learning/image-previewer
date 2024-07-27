package cache

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"io"
	"net/http"
	"strconv"

	"github.com/Roman-Otus-Learning/image-previewer/internal/app"
	"github.com/Roman-Otus-Learning/image-previewer/internal/cache/filesystem"
	"github.com/Roman-Otus-Learning/image-previewer/internal/cache/lru"
)

var _ app.App = (*AppCacheDecorator)(nil)

type AppCacheDecorator struct {
	app   app.App
	cache lru.Cache
	fs    filesystem.Filesystem
}

func CreateAppCacheDecorator(app app.App, limit uint64, cachePath string) (*AppCacheDecorator, error) {
	fs, err := filesystem.CreateDiskFilesystem(cachePath)
	if err != nil {
		return nil, fmt.Errorf("new cached app: %w", err)
	}

	return &AppCacheDecorator{
		app: app,
		cache: lru.CreateCacheLRU(limit, func(item *lru.Item) {
			_ = fs.RemoveFile(item.FileName)
		}),
		fs: fs,
	}, nil
}

func (a *AppCacheDecorator) Resize(
	ctx context.Context,
	url string,
	width, height int,
	headers http.Header,
) ([]byte, error) {
	cacheKey := a.generateCacheKey(url, width, height)

	item, found := a.cache.Get(cacheKey)
	if found {
		content, err := a.fs.ReadFile(item.FileName)
		if err != nil {
			return nil, errors.Wrap(err, "read from cached decorator")
		}

		log.Info().Msg("read file from from cache")

		return content, nil
	}

	content, err := a.app.Resize(ctx, url, width, height, headers)
	if err != nil {
		return nil, errors.Wrap(err, "resize from cached decorator")
	}

	item = &lru.Item{
		FileName: cacheKey + ".jpg",
		Size:     uint64(len(content)),
	}

	if err := a.fs.WriteFile(item.FileName, content); err != nil {
		return nil, errors.Wrap(err, "save in cached decorator")
	}
	a.cache.Set(cacheKey, item)

	return content, nil
}

func (a *AppCacheDecorator) generateCacheKey(url string, width, height int) string {
	hash := sha256.New()

	io.WriteString(hash, url)
	io.WriteString(hash, strconv.Itoa(width))
	io.WriteString(hash, strconv.Itoa(height))

	return hex.EncodeToString(hash.Sum(nil))
}
