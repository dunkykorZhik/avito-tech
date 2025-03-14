package main

import (
	"os"

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
	if err != nil || db == nil {
		logger.Fatalw("cannot connect to databse ", err)
		return
	}
	defer db.GetDb().Close()

	srv := server.NewHttpServer(cfg, logger, db)

	err = srv.ListenAndServe()
	logger.Error(err.Error())
	os.Exit(1)
}
