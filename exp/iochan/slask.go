
type IOChan struct {
	io.Closer
	in chan <- interface{}
	out <- chan interface{}
	quit chan bool
}

func (ioc IOChan) Close() (err error) {
	defer err = ioc.Closer.Close()
	ioc.quit <- true
	return
}

func NewReader(rc io.ReadCloser) (ioc *IOChan) {
	ioc = &IOChan{rc, make(chan interface{}), make(chan interface{}), make(chan bool)}
	return
}

func NewWriter(wc io.WriteCloser) (ioc *IOChan) {
	ioc = &IOChan{wc, make(chan interface{}), make(chan interface{}), make(chan bool)}
	return
}

