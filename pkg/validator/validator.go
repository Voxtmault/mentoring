package validator

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

func Init() *validator.Validate {
	validate := validator.New()
	validate.RegisterTagNameFunc(RegisterJSONTagNameFunc)
	validate.RegisterTagNameFunc(RegisterQueryTagNameFunc)

	return validate
}

func RegisterJSONTagNameFunc(fld reflect.StructField) string {
	name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
	if name == "-" {
		return ""
	}
	return name
}

func RegisterQueryTagNameFunc(fld reflect.StructField) string {
	name := strings.SplitN(fld.Tag.Get("query"), ",", 2)[0]
	if name == "-" {
		return ""
	}
	return name
}

func MangleErrorMessage(err error) []string {
	errArr := []string{}
	for _, item := range err.(validator.ValidationErrors) {
		errArr = append(errArr, fmt.Sprintf("%s: %s", item.Field(), item.Tag()))
	}

	return errArr
}
