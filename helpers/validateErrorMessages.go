package helpers

import (
	"github.com/go-playground/validator"
)

var validate *validator.Validate = validator.New()

type ValidationFieldError struct {
	Param   string
	Message string
}

type ValidationFieldErrors []ValidationFieldError

func (v *ValidationFieldErrors) ErrorMessage() any {
	return map[string]any{
		"error": *v,
	}
}

func (v *ValidationFieldError) ErrorMessage() any {
	return map[string]any{
		"error": *v,
	}
}

func msgForTag(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "This field is required"
	case "email":
		return "Invalid email"
	case "max":
		return "max is " + fe.Param()
	case "min":
		return "min is " + fe.Param()
	case "url":
		return "Invalid url" + fe.Param()
	}
	return ""
	// return fe.Error() // default error
}

func ValidateVar(field interface{}, paramName string, tag string) {
	if err := validate.Var(field, tag); err != nil {
		ve, _ := err.(validator.ValidationErrors)
		// return &ValidationFieldError{paramName, msgForTag(ve[0])}
		panic(ValidationFieldError{paramName, msgForTag(ve[0])})
	}
}

// func ValidateVar(field interface{}, paramName string, tag string) *ValidationFieldError {
// 	if err := validate.Var(field, tag); err != nil {
// 		ve, _ := err.(validator.ValidationErrors)
// 		// out[i] = ValidationFieldError{fe.StructField(), msgForTag(fe)}
// 		return &ValidationFieldError{paramName, msgForTag(ve[0])}
// 	}
// 	return nil
// }

// func ValidateStruct(fieldAndTags interface{}) ValidationFieldErrors {
// 	if err := validator.New().Struct(fieldAndTags); err != nil {
// 		ve, _ := err.(validator.ValidationErrors)
// 		out := make([]ValidationFieldError, len(ve))
// 		for i, fe := range ve {
// 			out[i] = ValidationFieldError{fe.StructField(), msgForTag(fe)}
// 		}
// 		return out
// 	}
// 	return nil
// }
