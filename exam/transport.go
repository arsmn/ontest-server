package exam

import (
	v "github.com/go-ozzo/ozzo-validation/v4"
)

///// CreateExamRequest

type CreateExamRequest struct {
	Title string `json:"title,omitempty"`
}

func (r CreateExamRequest) Validate() error {
	return v.ValidateStruct(&r,
		v.Field(&r.Title, v.Required, v.Length(5, 50)),
	)
}
