package iochan/codec

import (
	"iochan"
	"encoding/json"
)

type JsonCodec struct {
}

func NewJsonCodec() (*JsonCodec) {
	return &JsonCodec{}
}

func (JsonCodec) Marshal(i interface{}) ([]byte, error) {
	return json.Marshal(i)
}

func (JsonCodec) Unmarshal(b []byte, i interface{}) (error) {
	return json.Unmarshal(b, i)
}
