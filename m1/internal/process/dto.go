package process

import "github.com/google/uuid"

type RegisterDTO struct {
	TraceID         uuid.UUID
	Payload         []byte
	ByteSize        int
	TotalCharacters int
}

type ResponseRegisterDTO struct {
	TraceID string `json:"trace_id"`
}
