package transformers

import (
	"first/controllers/exceptions"

	"github.com/gofiber/fiber/v2"
)

type IResponseStruct struct {
	Status int64 `json:"status"`
	Data any `json:"data"`
}

func GlobalSuccessResponse(ctx *fiber.Ctx, res any) error {
	response := IResponseStruct{Status: 200, Data: res}
	err := ctx.JSON(response)
	if err != nil {
		return exceptions.ThrowInternalServerError(ctx)
	}
	return nil
}