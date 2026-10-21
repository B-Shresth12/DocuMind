package httpx

import (
	"errors"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func init() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			if name == "-" || name == "" {
				return fld.Name
			}
			return name
		})
	}
}

func FormatValidationErrors(err error) map[string]string {
	out := make(map[string]string)

	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		for _, fe := range ve {
			out[fe.Field()] = messageFor(fe)
		}

		return out
	}

	out["body"] = "invalid or malformed request body"
	return out
}

func messageFor(fe validator.FieldError) string {
	param := fe.Param()

	switch fe.Tag() {

	// Presence
	case "required":
		return "this field is required"
	case "required_with", "required_with_all", "required_without", "required_without_all", "required_if", "required_unless":
		return "this field is required"

	// Size / length (wording depends on the field's type)
	case "min", "gte":
		return "must be at least " + sizeText(fe.Kind(), param)
	case "max", "lte":
		return "must be at most " + sizeText(fe.Kind(), param)
	case "gt":
		return "must be greater than " + sizeText(fe.Kind(), param)
	case "lt":
		return "must be less than " + sizeText(fe.Kind(), param)
	case "len":
		return "must be exactly " + sizeText(fe.Kind(), param)

	// Equality / choices
	case "eq":
		return "must be equal to " + param
	case "ne":
		return "must not be equal to " + param
	case "oneof":
		return "must be one of: " + strings.ReplaceAll(param, " ", ", ")
	case "eqfield":
		return "must match " + param
	case "nefield":
		return "must be different from " + param

	// Formats
	case "email":
		return "must be a valid email address"
	case "url", "uri", "http_url":
		return "must be a valid URL"
	case "uuid", "uuid4":
		return "must be a valid UUID"
	case "alpha":
		return "must contain only letters"
	case "alphanum":
		return "must contain only letters and numbers"
	case "numeric", "number":
		return "must be a number"
	case "lowercase":
		return "must be lowercase"
	case "uppercase":
		return "must be uppercase"
	case "e164":
		return "must be a valid phone number in international format"
	case "datetime":
		return "must be a valid date/time in the format " + param
	case "timezone":
		return "must be a valid timezone"

	// String content
	case "contains":
		return "must contain " + param
	case "excludes":
		return "must not contain " + param
	case "startswith":
		return "must start with " + param
	case "endswith":
		return "must end with " + param

	// Collections
	case "unique":
		return "must not contain duplicate values"

	default:
		return "is invalid"
	}
}

// sizeText words a size constraint according to the field's kind:
// strings count characters, numbers are plain values, collections count items.
func sizeText(kind reflect.Kind, param string) string {
	switch kind {
	case reflect.String:
		return param + " characters"
	case reflect.Slice, reflect.Array, reflect.Map:
		return param + " items"
	default:
		return param
	}
}
