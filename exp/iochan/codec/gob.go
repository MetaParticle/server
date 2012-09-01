package iochan/codec

import (
	"iochan"
	"encoding/gob"
)

type GobCodec struct {
	enc *gob.Encoder
	dec *gob.Decoder
	buf bytes.Buffer
}

func NewGobCodec() (g *GobCodec) {
	g = &GobCodec{}
	g.enc = gob.NewEncoder(&g.buf)
	g.dec = gob.NewDecoder(&g.buf)
	return
}

func (g GobCodec) Marshal(i interface{}) (b []byte, err error) {
	err = g.enc.Encode(i)
	b = g.buf.Bytes()
	g.buf.Reset()
	return
}

func (g GobCodec) Unmarshal(b []byte, i interface{}) (err error) {
	g.buf.Write(b)
	/*
	switch reflect.TypeOf(i).(type) {
		case reflect.Value: {
			err = g.dec.DecodeValue(i)
		}
		default: {
			err = g.dec.Decode(i)
		}
	}
	*/
	err = g.dec.Decode(i)
	return
}

