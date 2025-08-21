package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dunkykorZhik/avito-tech/internal/app/config"
	"github.com/dunkykorZhik/avito-tech/internal/app/db"
	"github.com/dunkykorZhik/avito-tech/internal/app/server"
	"go.uber.org/zap"
)

func main() {
	logger := zap.Must(zap.NewProduction()).Sugar()
	defer logger.Sync()

	logger.Infow("Initializing configurations...")
	cfg, err := config.GetConfig()
	if err != nil {
		logger.Fatalw("cannot get config ", err)
		return
	}

	db, err := db.NewPostgresDB(cfg)
	if err != nil || db == nil || db.GetDb() == nil {
		logger.Fatalw("cannot connect to database ", err)
		return
	}
	defer db.GetDb().Close()

	srv := server.NewHttpServer(cfg, logger, db)
	shutdown := make(chan error)

	go func() {
		quit := make(chan os.Signal, 1)

		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		s := <-quit

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		logger.Infow("signal caught", "signal", s.String())

		shutdown <- srv.Shutdown(ctx)
	}()

	logger.Infof("Listening to server on port %d", cfg.Server.Port)
	err = srv.ListenAndServe()

	if !errors.Is(err, http.ErrServerClosed) {
		logger.Fatal(err)
	}

	err = <-shutdown
	if err != nil {
		logger.Fatal(err)
	}

}

// cors?
//metrics
