package network

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"strings"
)

func processErrors[T any](dto Dto[T], err error) error {
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		msgs, e := dto.ValidateErrors(validationErrors)
		if e != nil {
			return e
		}

		var msg strings.Builder
		br := ", "
		for _, m := range msgs {
			msg.WriteString(m + br)
		}
		// Remove the trailing separator
		errorMsg := msg.String()
		if len(errorMsg) > 0 {
			errorMsg = errorMsg[:len(errorMsg)-len(br)]
		}
		return errors.New(errorMsg)
	}

	return err
}

func ReqBody[T any](ctx *gin.Context, dto Dto[T]) (*T, error) {
	if err := ctx.ShouldBindJSON(dto); err != nil {
		e := processErrors(dto, err)
		return nil, e
	}

	v := validator.New()
	v.RegisterTagNameFunc(CustomTagNameFunc())
	if err := v.Struct(dto); err != nil {
		e := processErrors(dto, err)
		return nil, e
	}

	return dto.GetValue(), nil
}
