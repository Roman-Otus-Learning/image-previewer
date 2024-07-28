package builder

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog"
)

type shutdownFn func(context.Context) error

type shutdown struct {
	fn []shutdownFn
}

func (b *Builder) WaitShutdown(ctx context.Context) {
	stopSignals := []os.Signal{syscall.SIGTERM, syscall.SIGINT}
	s := make(chan os.Signal, len(stopSignals))
	signal.Notify(s, stopSignals...)
	zerolog.Ctx(ctx).Info().Msgf("got %s os signal. application will be stopped", <-s)

	b.shutdown.do(ctx)
}

func (s *shutdown) add(fn shutdownFn) {
	s.fn = append(s.fn, fn)
}

func (s *shutdown) do(ctx context.Context) {
	for i := len(s.fn) - 1; i >= 0; i-- {
		if err := s.fn[i](ctx); err != nil {
			zerolog.Ctx(ctx).Err(err).Send()
		}
	}
}
