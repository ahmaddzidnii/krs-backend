package utils

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"
)

type Response struct {
	Status           int               `json:"status"`
	ValidationErrors map[string]string `json:"validation_errors"`
	Message          string            `json:"message"`
	Data             interface{}       `json:"data"`
}

func Success(c *fiber.Ctx, statusCode int, message string, data interface{}) error {
	return c.Status(statusCode).JSON(Response{
		Status:           statusCode,
		Message:          message,
		ValidationErrors: nil,
		Data:             data,
	})
}

func Error(c *fiber.Ctx, statusCode int, message string) error {
	return c.Status(statusCode).JSON(Response{
		Status:           statusCode,
		Message:          message,
		ValidationErrors: nil,
		Data:             nil,
	})
}

func ValidationErrorResponse(c *fiber.Ctx, errors map[string]string) error {
	return c.Status(fiber.StatusBadRequest).JSON(Response{
		Status:           fiber.StatusBadRequest,
		Message:          "Validation failed. Please check the provided data.",
		ValidationErrors: errors,
		Data:             nil,
	})
}

func ClearCookies(c *fiber.Ctx, key ...string) {
	for i := range key {
		c.Cookie(&fiber.Cookie{
			Name:    key[i],
			Expires: time.Now().Add(-time.Hour * 24),
			Value:   "",
		})
	}
}

func GetLocals[T any](c *fiber.Ctx, key string) (T, error) {
	var result T // Variabel penampung dengan tipe generik T

	data := c.Locals(key)
	if data == nil {
		return result, fmt.Errorf("data dengan key '%s' tidak ditemukan injector context", key)
	}

	result, ok := data.(T)
	if !ok {
		// %T akan mencetak tipe data, berguna untuk debugging
		return result, fmt.Errorf("gagal konversi data: tipe data injector context (%T) tidak cocok dengan tipe yang diminta (%T)", data, result)
	}

	return result, nil
}

func GenerateValidationResponse(err error) map[string]string {
	errorBag := make(map[string]string)

	// Safely perform type assertion to validator.ValidationErrors
	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		// If the error is not from the validator, return a generic error message.
		// This prevents a panic if a different error type is passed.
		errorBag["error"] = "An unexpected error occurred during validation."
		return errorBag
	}

	for _, e := range validationErrors {
		// e.Field()  -> Struct field name (e.g., "Username")
		// e.Tag()    -> The failed rule (e.g., "required", "min")
		// e.Param()  -> The rule's parameter (e.g., "8" for "min=8")

		fieldName := strings.ToLower(e.Field())

		switch e.Tag() {
		case "required":
			errorBag[fieldName] = fmt.Sprintf("The %s field is required.", fieldName)
		case "email":
			errorBag[fieldName] = fmt.Sprintf("The %s field must be a valid email address.", fieldName)
		case "min":
			errorBag[fieldName] = fmt.Sprintf("The %s field must be at least %s characters long.", fieldName, e.Param())
		case "max":
			errorBag[fieldName] = fmt.Sprintf("The %s field may not be greater than %s characters.", fieldName, e.Param())
		case "len":
			errorBag[fieldName] = fmt.Sprintf("The %s field must be %s characters long.", fieldName, e.Param())
		case "eq":
			errorBag[fieldName] = fmt.Sprintf("The %s field must be equal to %s.", fieldName, e.Param())
		case "ne":
			errorBag[fieldName] = fmt.Sprintf("The %s field must not be equal to %s.", fieldName, e.Param())
		case "gt":
			errorBag[fieldName] = fmt.Sprintf("The %s field must be greater than %s.", fieldName, e.Param())
		case "gte":
			errorBag[fieldName] = fmt.Sprintf("The %s field must be greater than or equal to %s.", fieldName, e.Param())
		case "lt":
			errorBag[fieldName] = fmt.Sprintf("The %s field must be less than %s.", fieldName, e.Param())
		case "lte":
			errorBag[fieldName] = fmt.Sprintf("The %s field must be less than or equal to %s.", fieldName, e.Param())
		case "oneof":
			errorBag[fieldName] = fmt.Sprintf("The %s field must be one of the following: %s.", fieldName, strings.Replace(e.Param(), " ", ", ", -1))
		case "unique":
			errorBag[fieldName] = fmt.Sprintf("The %s field must contain unique values.", fieldName)
		case "url":
			errorBag[fieldName] = "The %s field must be a valid URL."
		case "uuid", "uuid4":
			errorBag[fieldName] = fmt.Sprintf("The %s field must be a valid UUID.", fieldName)
		case "alpha":
			errorBag[fieldName] = fmt.Sprintf("The %s field may only contain alphabetic characters.", fieldName)
		case "alphanum":
			errorBag[fieldName] = fmt.Sprintf("The %s field may only contain alpha-numeric characters.", fieldName)
		case "numeric":
			errorBag[fieldName] = fmt.Sprintf("The %s field must be a number.", fieldName)
		case "hostname":
			errorBag[fieldName] = fmt.Sprintf("The %s field must be a valid hostname.", fieldName)
		case "ip", "ipv4", "ipv6":
			errorBag[fieldName] = fmt.Sprintf("The %s field must be a valid IP address.", fieldName)
		case "cidr", "cidrv4", "cidrv6":
			errorBag[fieldName] = fmt.Sprintf("The %s field must be a valid CIDR notation.", fieldName)
		case "datetime":
			errorBag[fieldName] = fmt.Sprintf("The %s field does not match the format %s.", fieldName, e.Param())
		default:
			errorBag[fieldName] = fmt.Sprintf("The %s field is invalid.", fieldName)
		}
	}
	return errorBag
}

func ParseJsonBody[T any](c *fiber.Ctx) (T, error) {
	var result T

	if err := c.BodyParser(&result); err != nil {
		return result, fmt.Errorf("gagal mem-parsing body api: %w", err)
	}

	return result, nil
}

func FirstToLower(s string) string {
	r, size := utf8.DecodeRuneInString(s)
	if r == utf8.RuneError && size <= 1 {
		return s
	}
	lc := unicode.ToLower(r)
	if r == lc {
		return s
	}
	return string(lc) + s[size:]
}

func FirstToUpper(s string) string {
	r, size := utf8.DecodeRuneInString(s)
	if r == utf8.RuneError && size <= 1 {
		return s
	}
	uc := unicode.ToUpper(r)
	if r == uc {
		return s
	}
	return string(uc) + s[size:]
}
