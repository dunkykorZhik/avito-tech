package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/dunkykorZhik/avito-tech/internal/app/config"
	"github.com/dunkykorZhik/avito-tech/internal/app/db"
	"github.com/dunkykorZhik/avito-tech/internal/repo"
	"github.com/dunkykorZhik/avito-tech/internal/routes"
	"github.com/dunkykorZhik/avito-tech/internal/service"
	"go.uber.org/zap"
)

func NewHttpServer(cfg *config.Config, logger *zap.SugaredLogger, db db.Database) *http.Server {
	mux := http.NewServeMux()

	deps := service.ServiceDependencies{
		Repo:     repo.NewRepositories(db),
		SignKey:  cfg.Jwt.Signkey,
		TokenTTL: cfg.Jwt.TokenTTL,
	}
	services := service.NewService(deps)

	routes.AddRoutes(mux, services, logger)
	return &http.Server{
		//net.JoinHostPort(cfg.Server.Host, cfg.Server.Port),
		Addr:         fmt.Sprintf(cfg.Server.Port),
		Handler:      mux,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
}
