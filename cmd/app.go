package main

import (
	"context"
	"io"
	"os"
	"sync"

	"github.com/Roman-Otus-Learning/image-previewer/internal/builder"
	"github.com/Roman-Otus-Learning/image-previewer/internal/config"
	"github.com/pkg/errors"
	"github.com/rotisserie/eris"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type Application struct {
	Cfg *config.Config
}

func CreateApplication() *Application {
	return &Application{}
}

func (a *Application) RunCommands(rootCmd *cobra.Command) error {
	rootCmd.AddCommand(a.defineServerCmd())

	return rootCmd.Execute()
}

func (a *Application) defineServerCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "image-previewer",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			cfgPath, err := cmd.Flags().GetString("config")
			if err != nil {
				return errors.Wrap(err, "define config file path")
			}

			ctx, err = a.init(ctx, cfgPath)
			if err != nil {
				return err
			}

			logger := zerolog.Ctx(ctx)

			logger.Info().Msg("starting image previewer")
			appBuilder := builder.CreateBuilder(a.Cfg)

			clientInstance := appBuilder.CreateHTTPClient()
			resizerInstance := appBuilder.CreateResizer()
			application, err := appBuilder.CreateApplication(clientInstance, resizerInstance)
			if err != nil {
				return err
			}

			server := appBuilder.CreateHTTPServer(application)

			wg := &sync.WaitGroup{}
			wg.Add(1)
			server.Start(wg)
			appBuilder.WaitShutdown(ctx)
			wg.Wait()

			logger.Info().Msg("image previewer stopped")

			return nil
		},
	}

	return cmd
}

func (a *Application) init(ctx context.Context, cfgFilePath string) (context.Context, error) {
	var err error

	if a.Cfg, err = a.loadConfig(cfgFilePath); err != nil {
		return nil, errors.Wrap(err, "load config")
	}

	ctx = a.initLogger(ctx)

	return ctx, nil
}

func (a *Application) initLogger(ctx context.Context) context.Context {
	var writer io.Writer = os.Stdout

	logger := zerolog.New(writer).With().Timestamp().Logger().Level(zerolog.Level(a.Cfg.Logger.Level))
	log.Logger = logger

	ctxWithLogger := logger.WithContext(ctx)

	return ctxWithLogger
}

func (a *Application) loadConfig(cfgFilePath string) (*config.Config, error) {
	viper.SetConfigFile(cfgFilePath)
	if err := viper.ReadInConfig(); err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if errors.As(err, &configFileNotFoundError) {
			return nil, eris.New("Config file not found")
		}
	}

	cnf := &config.Config{}

	if err := viper.Unmarshal(cnf); err != nil {
		return nil, err
	}

	return cnf, nil
}
