package utils

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

//Overridding the validator to pick the json:"name" field for validation error instead of default Name
func NewValidator() *validator.Validate {
	validate := validator.New()
	validate.RegisterTagNameFunc(func(fl reflect.StructField) string {
		name := strings.SplitN(fl.Tag.Get("json"), ",", 2)
		if len(name) == 2 && name[1] == "-" {
			return ""
		}
		return name[0]
	})
	return validate
}

//this basically only prettify our errorFunc of the failed structs
func ValidationErrors(err error) map[string]string {
	errorsMap := make(map[string]string)

	if err == nil {
		return errorsMap
	}

	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		errorsMap["error"] = "Unknown error"
		return errorsMap
	}

	//we can also use fieldError.Param() to get the value let=10 i.e to get 10. For further beautify
	//the response messages 
	for _, fieldError := range validationErrors {
		msg := fmt.Sprintf(
			"failed '%s' tag check, value '%s' is not valid",
			fieldError.Tag(), fieldError.Value())
		errorsMap[fieldError.Field()] = msg
	}
	return errorsMap
}

// CheckForValidationError func for checking validation errors in struct fields.
func CheckForValidation(ctx *fiber.Ctx, errFunc error, statusCode int, object string) error {
	if errFunc != nil {
		return ctx.JSON(&fiber.Map{
			"status": statusCode,
			"msg":    fmt.Sprintf("validation errors for the %s fields", object),
			"fields": ValidationErrors(errFunc),
		})
	}
	return nil
}
