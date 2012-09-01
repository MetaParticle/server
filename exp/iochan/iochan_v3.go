package iochan

import (
	"reflect"
	"unsafe"
	"errors"
)

type Decoder interface {
	Decode(interface{}) error
}

type Encoder interface {
	Encode(interface{}) error
}

type inner func() error

func makeChan(typ reflect.Type, buffer int) (t_ch reflect.Value, ch chan <- interface{}) {
	t_ch := reflect.MakeChan(typ, buffer)
	p := t_ch.Pointer()
	ch = *p.(*chan interface{})
	return
}

func loop(inner) {
	var err error
	for err != nil {
		err = inner()
	}
}

func (enc Encoder) MakeChan(typ reflect.Type, buffer int) (ch chan <- interface{}) {
	t_ch, ch := makeChan(typ, buffer - 1)
	go loop(ch, func() (err error) {
		i, ok := t_ch.Recv()
		if ok {
			err = enc.Encode(i)
		} else {
			err = errors.New("Channel closed.")
		}
		return 
	})
	return
}

func (dec Decoder) MakeChan(typ reflect.Type, buffer int) (ch <- chan interface{}) {
	t_ch, ch := makeChan(typ, buffer - 1)
	go loop(func() (err error) {
		i := reflect.New(typ).Pointer()
		err = dec.Decode(i)
		t_ch.Send(i)
		// Bra att använda reflect's typ för kannaler och dess metod Send?
		return
	})
	return
}
