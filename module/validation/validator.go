package validation

import (
	ozzo "github.com/go-ozzo/ozzo-validation/v4"
)

func Validate(obj interface{}) error {
	if v, ok := obj.(ozzo.Validatable); ok {
		if err := v.Validate(); err != nil {
			if es, ok := err.(ozzo.Errors); ok {
				e := ErrValidation.WithDetail("code", "40010")
				for key, err := range es {
					if eo, ok := err.(ozzo.ErrorObject); ok {
						e.WithDetail(key, eo.Code())
					}
				}
				return e
			}
			return err
		}
	}
	return nil
}
