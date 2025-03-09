package fiberServer

import (
	"log"

	"github.com/dunkykorZhik/avito-tech/config"
	"github.com/dunkykorZhik/avito-tech/internal/db"
	"github.com/dunkykorZhik/avito-tech/internal/repo"
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

	f.app.Get("/hello", hello)

	log.Fatal(f.app.Listen(f.cfg.Server.Host + f.cfg.Server.Port))

}
func hello(c *fiber.Ctx) error {
	return c.SendString("Hello")

}

func GetUser(service service.User) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user, err := service.GetUser(c.Context(), "user1")
		log.Println(err)
		log.Println(user)

		return nil

	}

}

func (f *fiberServer) initializeHandlers() {
	deps := service.ServiceDependencies{
		Repo: repo.NewRepositories(f.db),
	}
	services := service.NewService(deps)
	f.app.Get("/user", GetUser(services.User))
}
