package utils

import (
	"reflect"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

func UserIdValidator() *validator.Validate {
	validate := validator.New()

	// Custom validation for UserId (positive integer)
	_ = validate.RegisterValidation("userId", func(fl validator.FieldLevel) bool {
		field := fl.Field()

		// Проверяем, что это целое число (int, int64 и т.д.)
		switch field.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return field.Int() > 0 // Должно быть > 0
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return true // Беззнаковые целые всегда >= 0
		case reflect.String:
			// Если пришла строка (например, из JSON), пытаемся преобразовать в int
			val, err := strconv.Atoi(field.String())
			if err != nil {
				return false
			}
			return val > 0
		default:
			return false // Все остальные типы недопустимы
		}
	})

	return validate
}

// NewValidator func for create a new validator for model fields.
func NewValidator() *validator.Validate {
	// Create a new validator for a Book model.
	validate := validator.New()

	// Custom validation for uuid.UUID fields.
	_ = validate.RegisterValidation("uuid", func(fl validator.FieldLevel) bool {
		field := fl.Field().String()
		if _, err := uuid.Parse(field); err != nil {
			return false // if there is an error, validation should return false
		}
		return true // if no error, validation should return true
	})

	return validate
}

// ValidatorErrors func for show validation errors for each invalid fields.
func ValidatorErrors(err error) map[string]string {
	// Define fields map.
	fields := map[string]string{}

	// Make error message for each invalid field.
	for _, err := range err.(validator.ValidationErrors) {
		fields[err.Field()] = err.Error()
	}

	return fields
}
