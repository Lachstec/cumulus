package api

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
)

var validate = validator.New()

func formatValidationErrors(errs validator.ValidationErrors) string {
	var errors []string
	for _, fieldErr := range errs {
		// fieldErr.Field() gives the name of the field that failed validation.
		// fieldErr.Tag() gives the validation rule that failed.
		errors = append(errors, fmt.Sprintf("Field '%s' failed on the '%s' tag", fieldErr.Field(), fieldErr.Tag()))
	}
	return strings.Join(errors, ", ")
}

func BindJSONStrict(c *gin.Context, obj interface{}) error {
	decoder := json.NewDecoder(c.Request.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(obj); err != nil {
		return err
	}

	if err := validate.Struct(obj); err != nil {
		// If the error is a ValidationErrors type, format the errors for better feedback.
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			return fmt.Errorf("validation failed: %s", formatValidationErrors(validationErrors))
		}
		return err
	}
	return nil
}
