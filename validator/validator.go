package validator

import (
	"reflect"
	"strings"
	"sync"

	"github.com/go-playground/validator/v10"
)

var (
	validate *validator.Validate
	once     sync.Once
)

func Init() {
	once.Do(func() {
		validate = validator.New()

		validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
			jsonTag := fld.Tag.Get("json")
			if jsonTag != "" {
				name := strings.SplitN(jsonTag, ",", 2)[0]
				if name != "-" && name != "" {
					return name
				}
			}

			formTag := fld.Tag.Get("form")
			if formTag != "" {
				name := strings.SplitN(formTag, ",", 2)[0]
				if name != "-" && name != "" {
					return name
				}
			}

			return fld.Name
		})

		registerCustomValidators(validate)
	})
}

func GetValidator() *validator.Validate {
	if validate == nil {
		Init()
	}
	return validate
}

func init() {
	Init()
}
