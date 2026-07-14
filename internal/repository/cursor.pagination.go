package repository

import (
	"encoding/base64"
	"encoding/json"
	"time"
)

type cursorPayload struct {
	Updated_at time.Time `json:"updated_at"`
	Id         string    `json:"id"`
}

// Encode cursor
func EncodeCursorPayload(updatedAt time.Time, id string) string {
	payload := cursorPayload{Updated_at: updatedAt, Id: id}
	b, _ := json.Marshal(&payload)
	return base64.StdEncoding.EncodeToString(b)
}

// Decode cursor
func decodeCursor(cursor string) (cursorPayload, error) {
	b, err := base64.StdEncoding.DecodeString(cursor)
	if err != nil {
		return cursorPayload{}, err
	}
	var payload cursorPayload
	json.Unmarshal(b, &payload)
	return payload, nil
}
