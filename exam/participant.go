package exam

import "time"

type Participant struct {
	UserID         uint64    `json:"user_id,omitempty"`
	ExamID         uint64    `json:"exam_id,omitempty"`
	ParticipatedAt time.Time `json:"participated_at,omitempty"`
}
