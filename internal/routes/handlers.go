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
		user := entity.User{
			ID: 1,
		}
		result, err := historyService.GetHistory(c.Context(), user.ID)
		if err != nil {
			return handleError(c, err)
		}
		return c.Status(http.StatusOK).JSON(result)

	}
}

type sendCoinRequest struct {
	toUser string `json:"toUser" validate:"required,max=100"`
	amount int64  `json:"amount" validate:"required,gt=0"`
}

func SendCoin(service service.Transfer) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Set(fiber.HeaderContentType, fiber.MIMETextPlainCharsetUTF8)
		//TODO change the hardcoded to user from context
		user := entity.User{
			ID:       1,
			Username: ":)",
		}
		var sendCoinRequest sendCoinRequest

		if err := c.BodyParser(sendCoinRequest); err != nil {
			return errorResponse(c, http.StatusBadRequest, errs.ErrInvalidReq.Error())
		}
		if err := Validate.Struct(sendCoinRequest); err != nil {
			return errorResponse(c, http.StatusBadRequest, errs.ValidationError(err))
		}
		transfer := entity.Transfer{
			Sender:   user.Username,
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
		c.Set(fiber.HeaderContentType, fiber.MIMETextPlainCharsetUTF8)
		user := entity.User{
			ID: 1,
		}
		item := c.Query("item")
		if err := Validate.Var(item, "required,max=50"); err != nil {
			return errorResponse(c, http.StatusBadRequest, errs.ValidationError(err))

		} //TODO validations

		if err := service.BuyItem(c.Context(), user.ID, item); err != nil {
			return handleError(c, err)
		}
		return c.SendStatus(http.StatusOK)
	}
}
