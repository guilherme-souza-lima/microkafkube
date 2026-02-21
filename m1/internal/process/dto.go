package process

import (
	"time"
)

type RegisterDTO struct {
	TraceID         string
	Payload         []byte
	ByteSize        int
	TotalCharacters int
}

type ResponseRegisterDTO struct {
	TraceID string `json:"trace_id"`
}

type RegistrationDTO struct {
	TraceID         string    `db:"trace_id"`
	Payload         []byte    `db:"payload"`
	ByteSize        int       `db:"byte_size"`
	TotalCharacters int       `db:"total_characters"`
	Published       bool      `db:"published_to_queue"`
	CreatedAt       time.Time `db:"created_at"`
}
