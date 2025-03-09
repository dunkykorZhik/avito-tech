package routes

import (
	"github.com/dunkykorZhik/avito-tech/internal/service"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

/*
/info  historyService
/send  transferService
/buy{item}  inventoryService

middleware
/auth


//error types:
-not enough balance
-no such user
-no such item
*/

var Validate *validator.Validate

func init() {
	Validate = validator.New(validator.WithRequiredStructEnabled())
}
func ShopRoutes(app fiber.App, services service.Service) {
	app.Get("/info", GetInfo(services.History))
	app.Post("/sendCoin", SendCoin(services.Transfer))
	app.Post("/buy/{item}", BuyItem(services.Inventory))

}
