package helper

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

type ValidationOutput map[string]string

func msgForTag(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "This field is required"
	case "email":
		return "Invalid email"
	case "min":
		return "Minimum character " + fe.Param()
	case "number":
		return "Must be a number"
	}

	return fe.Error()
}

func Validation(payload interface{}) ValidationOutput {
	validate = validator.New()
	errs := validate.Struct(payload)
	if errs != nil {
		fmt.Println(errs.Error())

		var apiErrors = make(ValidationOutput)
		for _, err := range errs.(validator.ValidationErrors) {
			fmt.Println()
			apiErrors[strings.ToLower(err.Field())] = msgForTag(err)
		}
		return apiErrors
	}

	return nil
}
