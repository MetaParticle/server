package iochan_v1

import (
	"io"
	"reflect"
	"encoding/gob"
)

type Encoder interface {
	Encode(interface{}) []byte
}

type Decoder interface {
	Decode([]byte) interface{}
}

type Codec interface {
	Encoder
	Decoder
}

func PipeWriter(ch chan interface{}, wr io.Writer, e Encoder ) {
	for {
		i := <- ch
		b := e.Encode(i)
		wr.Write(b)
		//wr.Write(f(<- ch))
	}
}

func PipeReader(ch chan interface{}, wr io.Reader, d Decoder ) {
	b := make([]byte, reflect.TypeOf(ch).Elem().Size())
	var (
		err error
	)
	for {
		_, err = wr.Read(b)
		if err == nil {
			i := d.Decode(b)
			ch <- i
		}
	}
}

type IOchan struct {
	typ reflect.Type
	chi, cho chan interface{}
	enc Codec
	rw io.ReadWriter
}

func NewIOchan(chi, cho chan interface{}, enc Codec, rw io.ReadWriter) (ioc *IOchan) {
	ioc = &IOchan{reflect.TypeOf(cho), chi, cho, enc, rw}
	return
}

func (ioc IOchan) Run() {
	go PipeReader(ioc.cho, ioc.rw, ioc.enc)
	go PipeWriter(ioc.chi, ioc.rw, ioc.enc)
}

type GobEncoder struct {
	size int
	buf []byte
	enc *gob.Encoder
	dec *gob.Decoder
}

type packet struct {
	typ string
	data []byte
}

func NewGobEncoder(typ interface{}) (gob *GobEncoder) {
	size := reflect.TypeOf(typ).Size()
	gob = &GobEncoder{size: size, buf: make([]byte, size)}
	gob.enc = gob.NewEncoder(gob)
	gob.dec = gob.NewDecoder(gob)
	return
}

func (gob GobEncoder) Read(p []byte) (n int, err error) {
	bl := len(gob.buf)
	pl := len(p)
	
	if bl <= pl {
		err = nil
		n = bl
	} else {
		err = io.ErrShortBuffer
		n = pl
	}
	
	for i := 0; i < n; i++ {
		p[i] = gob.buf[i]
	}
	gob.buf = gob.buf[n:]
	return
}

func (gob GobEncoder) Write(p []byte) (int, error) {
	gob.buf = append(gob.buf, p...)
	return len(p), nil
}

func (gob GobEncoder) Encode(i interface{}) (b []byte) {
	gob.enc.Encode(i)
	p := packet(reflect.TypeOf(i).String(), gob.buf)
	gob.buf = make([]byte, gob.size)
	gob.enc.Encode(p)
	b = gob.buf
	gob.buf = make([]byte, gob.size)
	return
}

func (gob GobEncoder) Decode(b []byte) (i interface{}) {
	return nil
}

