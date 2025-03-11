package routes

import (
	"net/http"

	"github.com/dunkykorZhik/avito-tech/internal/entity"
	errs "github.com/dunkykorZhik/avito-tech/internal/errors"
	"github.com/dunkykorZhik/avito-tech/internal/service"
	"github.com/gofiber/fiber/v2"
)

func GetInfo(historyService service.History) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
		//TODO change the hardcoded to user from context
		username := ""
		result, err := historyService.GetHistory(c.Context(), username)
		if err != nil {
			return handleError(c, err)
		}
		return c.Status(fiber.StatusOK).JSON(result)

	}
}

type sendCoinRequest struct {
	toUser string `json:"toUser" validate:"required,max=100"`
	amount int64  `json:"amount" validate:"required,gt=0"`
}

func SendCoin(service service.Transfer) fiber.Handler {
	return func(c *fiber.Ctx) error {

		//TODO change the hardcoded to user from context
		username := ""
		var sendCoinRequest sendCoinRequest

		if err := c.BodyParser(&sendCoinRequest); err != nil {
			return errorResponse(c, fiber.StatusBadRequest, errs.ErrInvalidReq.Error())
		}
		if err := Validate.Struct(sendCoinRequest); err != nil {
			return errorResponse(c, fiber.StatusBadRequest, errs.ValidationError(err))
		}
		transfer := entity.Transfer{
			Sender:   username,
			Receiver: sendCoinRequest.toUser,
			Amount:   sendCoinRequest.amount,
		}

		err := service.CreateTransfer(c.Context(), transfer)
		if err != nil {
			return handleError(c, err)
		}
		return c.SendStatus(http.StatusOK)

	}
}

func BuyItem(service service.Inventory) fiber.Handler {
	return func(c *fiber.Ctx) error {
		//TODO change the hardcoded to user from context

		username := ""
		item := c.Query("item")
		if err := Validate.Var(item, "required,max=50"); err != nil {
			return errorResponse(c, fiber.StatusBadRequest, errs.ValidationError(err))

		} //TODO validations

		if err := service.BuyItem(c.Context(), username, item); err != nil {
			return handleError(c, err)
		}
		return c.SendStatus(fiber.StatusOK)
	}
}

type AuthRequest struct {
	username string `json:"username" validate:"required,max=100"`
	password string `json:"password" validate:"required,min=6,max=100"`
}

func Auth(service service.User) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Set(fiber.HeaderAccept, fiber.MIMEApplicationJSON)
		var authReq AuthRequest
		if err := c.BodyParser(&authReq); err != nil {
			return errorResponse(c, fiber.StatusBadRequest, errs.ErrInvalidReq.Error())
		}
		if err := Validate.Struct(authReq); err != nil {
			return errorResponse(c, fiber.StatusBadRequest, errs.ValidationError(err))
		}
		token, err := service.GenerateToken(c.Context(), authReq.username, authReq.password)
		if err != nil {
			return handleError(c, err)
		}
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"token": token})

	}

}
