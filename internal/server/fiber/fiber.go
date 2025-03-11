package fiberServer

import (
	"log"

	"github.com/dunkykorZhik/avito-tech/config"
	"github.com/dunkykorZhik/avito-tech/internal/db"
	"github.com/dunkykorZhik/avito-tech/internal/repo"
	"github.com/dunkykorZhik/avito-tech/internal/routes"
	"github.com/dunkykorZhik/avito-tech/internal/service"
	"github.com/gofiber/fiber/v2"
)

type fiberServer struct {
	app *fiber.App
	cfg *config.Config
	db  db.Database
}

func NewFiberServer(cfg *config.Config, db db.Database) *fiberServer {
	fiberApp := fiber.New()
	return &fiberServer{
		app: fiberApp,
		cfg: cfg,
		db:  db,
	}
}
func (f *fiberServer) Start() {
	f.initializeHandlers()

	log.Fatal(f.app.Listen(f.cfg.Server.Host + f.cfg.Server.Port))

}

func (f *fiberServer) initializeHandlers() {
	deps := service.ServiceDependencies{
		Repo: repo.NewRepositories(f.db),
	}
	services := service.NewService(deps)
	api := f.app.Group("/api")
	routes.ShopRoutes(api, *services)

}
