package validator

import (
	"fmt"
	"log/slog"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		// Optionally, you could return the error to give each route more control over the status code
		return err
	}
	return nil
}

var validate *validator.Validate

func InitValidator() error {
	validate = validator.New()
	validate.RegisterTagNameFunc(RegisterTagNameFunc)

	slog.Debug("successfully initialized validator")
	return nil
}

func GetValidator() *validator.Validate {
	return validate
}

func GetValidatorInstance() *CustomValidator {
	return &CustomValidator{validator: validate}
}

func RegisterTagNameFunc(fld reflect.StructField) string {
	// List of tags to check, in order of priority
	tags := []string{"json", "query"}

	for _, tag := range tags {
		name := strings.SplitN(fld.Tag.Get(tag), ",", 2)[0]
		if name != "" && name != "-" {
			return name
		}
	}

	return ""
}

func MangleErrorMessage(err error) []string {
	errArr := []string{}
	if validationErr, ok := err.(validator.ValidationErrors); ok {
		for _, item := range validationErr {
			errArr = append(errArr, fmt.Sprintf("%s: %s", item.Field(), item.Tag()))
		}
		return errArr
	} else {
		return []string{err.Error()}
	}
}
