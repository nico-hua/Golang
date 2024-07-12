package out

import "fmt"

type Out struct {
	data chan interface{}
}

var out *Out

func NewOut() *Out {
	if out == nil {
		out = &Out{
			data: make(chan interface{}, 3),
		}
	}
	return out
}

func Println(i interface{}) {
	out.data <- i
}

func (o *Out) OutPut() {
	for {
		select {
		case i := <-o.data:
			fmt.Println(i)
		}
	}
}