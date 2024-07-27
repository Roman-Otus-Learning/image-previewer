package builder

import (
	"code.cloudfoundry.org/bytefmt"
	"github.com/Roman-Otus-Learning/image-previewer/internal/app"
	"github.com/Roman-Otus-Learning/image-previewer/internal/cache"
	"github.com/Roman-Otus-Learning/image-previewer/internal/client"
	"github.com/Roman-Otus-Learning/image-previewer/internal/config"
	"github.com/Roman-Otus-Learning/image-previewer/internal/resizer"
	"github.com/Roman-Otus-Learning/image-previewer/internal/server"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"time"
)

type Builder struct {
	config   *config.Config
	shutdown shutdown
}

func CreateBuilder(config *config.Config) *Builder {
	return &Builder{config: config}
}

func (b *Builder) CreateHTTPServer(app app.App) *server.Server {
	HTTPServer := server.CreateHTTPServer(b.config.HTTPAddr(), app)

	b.shutdown.add(
		func(ctx context.Context) error {
			HTTPServer.Stop(ctx)

			return nil
		},
	)

	return HTTPServer
}

func (b *Builder) CreateHTTPClient() *client.HTTPClient {
	return client.CreateHTTPClient(time.Duration(b.config.Client.Timeout))
}

func (b *Builder) CreateResizer() *resizer.ImageResizer {
	return resizer.CreateResizer()
}

func (b *Builder) CreateApplication(client *client.HTTPClient, resizer *resizer.ImageResizer) (app.App, error) {
	application := app.CreateResizerApp(client, resizer)

	cacheSizeBytes, err := bytefmt.ToBytes(b.config.Cache.Size)
	if err != nil {
		return nil, errors.Wrap(err, "invalid cache size")
	}

	cachedApp, err := cache.CreateAppCacheDecorator(application, cacheSizeBytes, b.config.Cache.Path)
	if err != nil {
		return nil, errors.Wrap(err, "create cached app")
	}

	return cachedApp, nil
}
