package exam

import (
	"time"

	"github.com/arsmn/ontest-server/user"
	v "github.com/go-ozzo/ozzo-validation/v4"
)

///// CreateExamRequest

type CreateExamRequest struct {
	user.SignedRequest
	Title    string     `json:"title,omitempty"`
	StartAt  time.Time  `json:"start_at,omitempty"`
	Once     bool       `json:"once,omitempty"`
	Deadline *time.Time `json:"deadline,omitempty"`
}

func (r CreateExamRequest) Validate() error {
	return v.ValidateStruct(&r,
		v.Field(&r.Title, v.Required, v.Length(5, 50)),
		v.Field(&r.StartAt, v.Required, v.Min(time.Now().UTC())),
		v.Field(&r.Deadline, v.When(!r.Once, v.Required).Else(v.Nil), v.Min(r.StartAt)),
	)
}
