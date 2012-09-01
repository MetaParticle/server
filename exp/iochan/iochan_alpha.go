package iochan

import (
	"io"
	"bytes"
	"reflect"
)

func CopyChan(chi <- chan interface{}, cho chan <- interface{}) {
	for {
		cho <- <- chi
	}
}

func WriteChan(w io.Writer, ch <- chan []byte) {
	for {
		w.Write(<- ch)
	}
}

func ReadChan(r io.Reader, ch chan <- []byte) (err error) {
	for {
		buf := make([]byte, 0)
		_, err = r.Read(buf)
		if err == nil {
			ch <- buf
		} else {
			break
		}
	}
	return
}

type Marshal func(i interface{}) ([]byte, error)
type Unmarshal func(b []byte, i interface{}) (error)

type Codec interface {
	Marshal(i interface{}) ([]byte, error)
	Unmarshal(b []byte, i interface{}) (error)
}

func MarshalChan(m Marshal, chi <- chan interface{}, cho chan <- []byte) {
	for {
		buf, err := m(<- chi)
		if err == nil {
			cho <- buf
		} else {
			break
		}
	}
}

func UnmarshalChan(um Unmarshal, chi <- chan []byte, cho chan <- interface{}) {
	for {
		i := reflect.Indirect(reflect.New(reflect.TypeOf(cho).Elem()))
		err := um(<- chi, i)
		if err == nil {
			cho <- i
		}
	}
}

func WriteCodecChan(w io.Writer, c *Codec, chi <- chan interface{}) {
	cho := make(chan []byte)
	go WriteChan(w, cho)
	go MarshalChan(c.Marshal, chi, cho)
	go CopyChan(chi, cho)
}

func ReadCodecChan(r io.Reader, c *Codec, cho chan <- interface{}) {
	chi := make(chan []byte)
	go ReadChan(r, chi)
	go UnmarshalChan(c.Unmarshal, chi, cho)
	go CopyChan(chi, cho)
}

type packet struct {
	typ reflect.Type
	data []byte
}

type SafeCodec struct {
	typ reflect.Type
	codec Codec
}

func NewSafeCodec(t interface{}, codec *Codec) (*SafeCodec) {
	return &SafeCodec{reflect.TypeOf(t), *codec}
}

func (s SafeCodec) Marshal(i interface{}) (b []byte, err error) {
	if reflect.TypeOf(i) != s.typ {
		//Vad sägs om att retunera nil istället? Nu måste man krångla med recover ifall något går fel...
		panic("Not the of the same type!")
	}
	b, err = s.codec.Marshal(i)
	if err != nil {
		return
	}
	p := packet{s.typ, b}
	b, err = s.codec.Marshal(p)
	return
}

func (s SafeCodec) Unmarshal(b []byte, i interface{}) (err error) {
	if reflect.TypeOf(i) != s.typ {
		//se ovan ^
		panic("Not the of the same type!")
	}
	p := packet{}
	err = s.codec.Unmarshal(b, p)
	if err != nil {
		return
	}
	if p.typ != s.typ {
		//se ovan ^
		panic("Not the of the same type!")
	}
	err = s.codec.Unmarshal(p.data, i)
	return
}

func (s SafeCodec) SafeUnmarshal(b []byte) (i interface{}, err error) {
	
	p := packet{}
	err = s.codec.Unmarshal(b, p)
	if err != nil {
		return
	}
	if p.typ != s.typ {
		//se ovan ^
		panic("Not the of the same type!")
	}
	i = reflect.New(p.typ)
	err = s.codec.Unmarshal(p.data, i)
	return
}


