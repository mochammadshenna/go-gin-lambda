package json

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// Marshal marshals an interface to JSON bytes
func Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

// Unmarshal unmarshals JSON bytes to an interface
func Unmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

// MarshalIndent marshals an interface to indented JSON bytes
func MarshalIndent(v interface{}, prefix, indent string) ([]byte, error) {
	return json.MarshalIndent(v, prefix, indent)
}

// MarshalToString marshals an interface to JSON string
func MarshalToString(v interface{}) (string, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return "", fmt.Errorf("failed to marshal to string: %w", err)
	}
	return string(data), nil
}

// UnmarshalFromString unmarshals a JSON string to an interface
func UnmarshalFromString(data string, v interface{}) error {
	return json.Unmarshal([]byte(data), v)
}

// IsValid checks if a string is valid JSON
func IsValid(data string) bool {
	var js json.RawMessage
	return json.Unmarshal([]byte(data), &js) == nil
}

// PrettyPrint prints JSON in a pretty format
func PrettyPrint(v interface{}) (string, error) {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to pretty print: %w", err)
	}
	return string(data), nil
}

// NewDecoder creates a new JSON decoder from a bytes.Buffer
func NewDecoder(buffer *bytes.Buffer) *json.Decoder {
	return json.NewDecoder(buffer)
}
