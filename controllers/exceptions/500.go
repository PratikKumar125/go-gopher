package exceptions

import "github.com/gofiber/fiber/v2"

type IInternalServerResponse struct {
	Status int    `json:"status"`
	Error  string `json:"error"`
}

func ThrowInternalServerError(ctx *fiber.Ctx) error {
	response := IInternalServerResponse{
		Status: 500,
		Error: "Internal Server Error",
	}
	return ctx.Status(fiber.StatusNotFound).JSON(response)
}