package validators

import (
	"ai-service/internal/model/api"
	"ai-service/internal/util/exception"
	"ai-service/internal/util/exceptioncode"
	"reflect"
	"regexp"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
)

var Validator *validator.Validate

func New() *validator.Validate {
	Validator = validator.New()
	Validator.RegisterTagNameFunc(extractJsonTag)

	// validate data type
	exception.PanicOnError(Validator.RegisterValidation("datetime_rfc3339", datetimeRFC3339))
	exception.PanicOnError(Validator.RegisterValidation("date", date))
	exception.PanicOnError(Validator.RegisterValidation("date_v2", dateV2))
	exception.PanicOnError(Validator.RegisterValidation("sort", sortParams))
	exception.PanicOnError(Validator.RegisterValidation("sort_strings", sortStrings))
	exception.PanicOnError(Validator.RegisterValidation("date_range", dateRangeParams))
	exception.PanicOnError(Validator.RegisterValidation("datetime_range", datetimeRange))
	exception.PanicOnError(Validator.RegisterValidation("language_number_type", languageNumberType))

	return Validator
}

func extractJsonTag(fld reflect.StructField) string {
	name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]

	if name == "-" {
		return ""
	}

	return name
}

func Validate(e interface{}) error {
	err := Validator.Struct(e)
	if err == nil {
		return err
	}

	errors := processValidationErrors(err)

	return api.ErrorResponse{
		CodeMessage: exceptioncode.CodeInvalidValidation,
		Message:     "validation error : " + errors[0].Message,
		Errors:      errors,
	}
}

func processValidationErrors(err error) []api.ErrorValidate {
	if err == nil {
		return nil
	}

	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		return []api.ErrorValidate{{
			Key:         "unknown",
			CodeMessage: "VALIDATION",
			Message:     err.Error(),
		}}
	}

	errors := make([]api.ErrorValidate, len(validationErrors))
	for i, er := range validationErrors {
		errors[i] = api.ErrorValidate{
			Key:         er.Field(),
			CodeMessage: "VALIDATION",
			Message:     er.Error(),
		}
	}

	return errors
}

func datetimeRFC3339(fl validator.FieldLevel) bool {
	datetime := fl.Field().String()

	_, err := time.Parse(time.RFC3339, datetime)
	return err == nil
}

func date(fl validator.FieldLevel) bool {
	datetime := fl.Field().String()

	_, err := time.Parse("2006-01-02", datetime)
	return err == nil
}

func dateV2(fl validator.FieldLevel) bool {
	datetime := fl.Field().String()

	t, err := time.Parse("2006-01-02", datetime)
	return err == nil && t.Year() >= 1900
}

func dateRangeParams(fl validator.FieldLevel) bool {
	s := fl.Field().String()

	return !(string(s[10]) != "~")
}

func datetimeRange(fl validator.FieldLevel) bool {
	s := fl.Field().String()
	split := strings.Split(s, "~")

	_, errI := time.Parse(time.RFC3339, split[0])
	_, errJ := time.Parse(time.RFC3339, split[1])
	if errI != nil || errJ != nil {
		return false
	}
	return true
}

func sortParams(fl validator.FieldLevel) bool {
	s := fl.Field().String()
	split := strings.Split(s, ",")
	for _, v := range split {
		if string(v[0]) != "+" && string(v[0]) != "-" {
			return false
		}
	}
	return true
}

// support only "+key" and "-key"
func sortStrings(fl validator.FieldLevel) bool {
	s := fl.Field().Interface().([]string)
	r, _ := regexp.Compile(`^[+\-][a-z_]+$`)
	for _, v := range s {
		if !r.MatchString(v) {
			return false
		}
	}
	return true
}

func languageNumberType(fl validator.FieldLevel) bool {
	f := fl.Field()

	switch f.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return f.Int() >= 1
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return f.Uint() >= 1
	case reflect.Float32, reflect.Float64:
		return f.Float() >= 1
	}

	// For custom numeric types
	if f.CanInterface() {
		switch v := f.Interface().(type) {
		case int:
			return v >= 1
		case int8:
			return v >= 1
		case int16:
			return v >= 1
		case int32:
			return v >= 1
		case int64:
			return v >= 1
		case uint:
			return v >= 1
		case uint8:
			return v >= 1
		case uint16:
			return v >= 1
		case uint32:
			return v >= 1
		case uint64:
			return v >= 1
		case float32:
			return v >= 1
		case float64:
			return v >= 1
		}
	}

	return false
}
