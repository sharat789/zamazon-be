package rest

import (
	"github.com/gofiber/fiber/v2"
	"net/http"
)

func ErrorResponse(ctx *fiber.Ctx, status int, err error) error {
	return ctx.Status(status).JSON(err.Error())
}

func InternalErrorResponse(ctx *fiber.Ctx, err error) error {
	return ctx.Status(http.StatusInternalServerError).JSON(err.Error())
}

func BadRequestErrorResponse(ctx *fiber.Ctx, msg string) error {
	return ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{
		"message": msg,
	})
}
func SuccessResponse(ctx *fiber.Ctx, message string, data interface{}) error {
	return ctx.Status(http.StatusOK).JSON(&fiber.Map{
		"message": message,
		"data":    data,
	})
}
