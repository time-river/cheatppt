package revchatgpt2

import (
	"encoding/json"
)

type marshaller interface {
	marshal(value any) ([]byte, error)
}

type jsonMarshaller struct{}

func (jm *jsonMarshaller) marshal(value any) ([]byte, error) {
	return json.Marshal(value)
}

type unmarshaler interface {
	unmarshal(data []byte, v any) error
}

type jsonUnmarshaler struct{}

func (jm *jsonUnmarshaler) unmarshal(data []byte, v any) error {
	return json.Unmarshal(data, v)
}
