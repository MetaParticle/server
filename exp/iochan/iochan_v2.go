package iochan

import (
	"reflect"
	"unsafe"
)

type Decoder interface {
	Decode(interface{}) error
}

type Encoder interface {
	Encode(interface{}) error
}

func NewInChan(enc *Encoder, typ interface{}, buffer int) (ch chan <- interface{}) {
	//ch = make(chan interface{}) //
	t_ch := &reflect.MakeChan(reflect.Type(typ), buffer - 1)
	p := t_ch.(unsafe.Pointer)
	ch = *p.(*chan interface{})
	go func() {
		var err error
		for err != nil {
			i := <- ch
			err = enc.Encode(i)
		}
	}()
	return
}

func NewOutChan(dec *Decoder, typ interface{}, buffer int) (ch <- chan interface{}) {
	t_ch := &reflect.MakeChan(reflect.Type(typ), buffer - 1)
	p := t_ch.(unsafe.Pointer)
	ch = *p.(*chan interface{})
	go func() {
		var err error
		for err != nil {
			i := &reflect.New(reflect.Type(typ))
			err = dec.Decode(i)
			ch <- i.(unsafe.Pointer)
		}
	}()
}
