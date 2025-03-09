package routes

import (
	errs "github.com/dunkykorZhik/avito-tech/internal/errors"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

func handleError(c *fiber.Ctx, err error) error {

	var data string
	status := fiber.StatusBadRequest
	if err == errs.ErrNoItem || err == errs.ErrNoUser || err == errs.ErrNotEnoughBalance {
		data = err.Error()
	} else {
		status = fiber.StatusInternalServerError
		log.Errorw(err.Error())
		data = "internal server error, try again later"
	}
	return errorResponse(c, status, data)

}

func errorResponse(c *fiber.Ctx, status int, data string) error {
	c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
	return c.Status(status).JSON(fiber.Map{
		"data": data,
	})

}
